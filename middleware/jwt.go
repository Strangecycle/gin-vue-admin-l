package middleware

import (
	"errors"
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model"
	"gin-vue-admin-l/model/request"
	"gin-vue-admin-l/model/response"
	"gin-vue-admin-l/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"time"
)

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

type JWT struct {
	SigningKey []byte
}

func NewJWT() *JWT {
	return &JWT{[]byte(global.GVA_CONFIG.JWT.SigningKey)}
}

// 创建 token
func (j *JWT) CreateToken(c request.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(j.SigningKey)
}

// 解析 token
func (j *JWT) ParseToken(tStr string) (*request.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tStr, &request.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		// 判断 token 无效的原因
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed // token 格式错误
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired // token 过期
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet // token 暂未验证
			} else {
				return nil, TokenInvalid // token 无效
			}
		}
	}

	if token != nil {
		if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}

// JWT 鉴权中间件
// 返回 reload: true 表示鉴权失败q前端需要刷新重定向到登录页重新登录
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("x-token")
		// 前端未传递 token
		if token == "" {
			response.FailWithDetailed(gin.H{"reload": true}, "未登录或非法访问", c)
			c.Abort()
			return
		}

		// 当 token 在黑名单中时，说明已经在别处登录过，此时通知前端重定向，从而实现了单点登录的功能
		if service.IsInBlacklist(token) {
			response.FailWithDetailed(gin.H{"reload": true}, "您的帐户异地登陆或令牌失效", c)
			c.Abort()
			return
		}

		// 开始解析 token
		j := NewJWT()
		claims, err := j.ParseToken(token)
		// 解析 token 失败
		if err != nil {
			if err == TokenExpired {
				response.FailWithDetailed(gin.H{"reload": true}, "授权已过期", c)
				c.Abort()
				return
			}
			response.FailWithDetailed(gin.H{"reload": true}, err.Error(), c)
			c.Abort()
			return
		}

		// token 解析成功,查找用户
		if err, _ := service.FindUserByUuid(claims.UUID); err != nil {
			// 未找到用户,将 token 拉黑
			_ = service.JsonInBlacklist(model.JwtBlacklist{Jwt: token})
			response.FailWithDetailed(gin.H{"reload": true}, err.Error(), c)
			c.Abort()
			return
		}

		// 如果当前的 token 剩余时间小于缓冲时间,则代表 token 即将过期,进入自动刷新 token 流程
		// 用于确保 token 不会出现过期情况
		if claims.ExpiresAt-time.Now().Unix() < claims.BufferTime {
			// 重新设置 token 过期时间
			claims.ExpiresAt = time.Now().Unix() + global.GVA_CONFIG.JWT.ExpiresTime
			// 生成新的 token 与 claims
			newToken, _ := j.CreateToken(*claims)
			newClaims, _ := j.ParseToken(newToken)
			c.Header("new-token", newToken)
			c.Header("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt, 10))
			// 如果是单点登录模式,则刷新 redis 中的 token
			if global.GVA_CONFIG.System.UseMultipoint {
				err, redisJWT := service.GetRedisJWT(newClaims.Username)
				if err != nil {
					global.GVA_LOG.Error("get redis jwt failed", zap.Any("err", err))
				} else {
					_ = service.JsonInBlacklist(model.JwtBlacklist{Jwt: redisJWT})
				}
				service.SetRedisJWT(newToken, newClaims.Username)
			}
		}
		c.Set("claims", claims)
		c.Next()
	}
}
