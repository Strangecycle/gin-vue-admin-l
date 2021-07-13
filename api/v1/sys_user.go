package v1

import (
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/middleware"
	"gin-vue-admin-l/model"
	"gin-vue-admin-l/model/request"
	"gin-vue-admin-l/model/response"
	"gin-vue-admin-l/service"
	"gin-vue-admin-l/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"time"
)

// 签发 token
func tokenNext(c *gin.Context, user model.SysUser) {
	j := &middleware.JWT{SigningKey: []byte(global.GVA_CONFIG.JWT.SigningKey)}
	claims := request.CustomClaims{
		UUID:        user.UUID,
		ID:          user.ID,
		NickName:    user.NickName,
		AuthorityId: user.AuthorityId,
		BufferTime:  global.GVA_CONFIG.JWT.BufferTime, // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,                              // 签名生效时间
			ExpiresAt: time.Now().Unix() + global.GVA_CONFIG.JWT.ExpiresTime, // 过期时间 7天  配置文件
			Issuer:    "qmPlus",                                              // 签名的发行者
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		global.GVA_LOG.Error("获取token失败!", zap.Any("err", err))
		response.FailWithMessage("获取token失败", c)
		return
	}

	if !global.GVA_CONFIG.System.UseMultipoint {
		response.OkWithDetailed(response.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功", c)
	}

	// redis 多点登录拦截 - 单点登录
	if err, jwtStr := service.GetRedisJWT(user.Username); err == redis.Nil {
		// redis 中不存在数据，没有登录过
		if err := service.SetRedisJWT(token, user.Username); err != nil {
			global.GVA_LOG.Error("设置登录状态失败!", zap.Any("err", err))
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		response.OkWithDetailed(response.LoginResponse{
			User:      user,
			Token:     jwtStr,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功", c)
	} else if err != nil {
		global.GVA_LOG.Error("设置登录状态失败!", zap.Any("err", err))
		response.FailWithMessage("设置登录状态失败", c)
	} else {
		// 已经登陆过，第二次登录或异地登录 - 将其他人挤下线
		var blackJWT model.JwtBlacklist
		blackJWT.Jwt = jwtStr
		// 第二次登录时将当前的 token 加入黑名单作废
		if err := service.JsonInBlacklist(blackJWT); err != nil {
			response.FailWithMessage("jwt作废失败", c)
			return
		}
		// 重新赋值 redis 中的 token，表示是一个新的登录
		if err := service.SetRedisJWT(token, user.Username); err != nil {
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		response.OkWithDetailed(response.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功", c)
	}
}

// @Tags Base
// @Summary 用户登录
// @Produce  application/json
// @Param data body request.Login true "用户名, 密码, 验证码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登陆成功"}"
// @Router /base/login [post]
func Login(c *gin.Context) {
	var l request.Login
	_ = c.ShouldBindJSON(&l)
	// 验证登录参数
	if err := utils.Verify(l, utils.LoginVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 验证码库校验方法
	if !store.Verify(l.CaptchaId, l.Captcha, true) {
		response.FailWithMessage("验证码错误", c)
		return
	}

	u := model.SysUser{Username: l.Username, Password: l.Password}
	err, user := service.Login(&u)
	if err != nil {
		global.GVA_LOG.Error("登陆失败! 用户名不存在或者密码错误!", zap.Any("err", err))
		response.FailWithMessage("用户名不存在或者密码错误", c)
		return
	}
	tokenNext(c, *user)
}

// @Tags SysUser
// @Summary 分页获取用户列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /user/getUserList [post]
func GetUserList(c *gin.Context) {
	var pageInfo request.PageInfo
	_ = c.ShouldBindJSON(&pageInfo)
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err, list, total := service.GetUserInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// @Tags SysUser
// @Summary 用户注册账号
// @Produce  application/json
// @Param data body model.SysUser true "用户名, 昵称, 密码, 角色ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"注册成功"}"
// @Router /user/register [post]
func Register(c *gin.Context) {
	var r request.Register
	_ = c.ShouldBindJSON(&r)
	if err := utils.Verify(r, utils.RegisterVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	u := &model.SysUser{
		Username:    r.Username,
		Password:    r.Password,
		NickName:    r.NickName,
		HeaderImg:   r.NickName,
		AuthorityId: r.AuthorityId,
	}
	err, user := service.Register(*u)
	if err != nil {
		global.GVA_LOG.Error("注册失败!", zap.Any("err", err))
		response.FailWithDetailed(response.SysUserResponse{User: user}, "注册失败", c)
		return
	}
	response.OkWithDetailed(response.SysUserResponse{User: user}, "注册成功", c)
}

// @Tags SysUser
// @Summary 删除用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "用户ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /user/deleteUser [delete]
func DeleteUser(c *gin.Context) {
	var reqId request.GetById
	_ = c.ShouldBindJSON(&reqId)
	if err := utils.Verify(reqId, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	jwtId := getUserID(c)
	if jwtId == uint(reqId.ID) {
		response.FailWithMessage("删除失败, 自杀失败", c)
		return
	}

	err := service.DeleteUser(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// @Tags SysUser
// @Summary 设置用户权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SetUserAuth true "用户UUID, 角色ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /user/setUserAuthority [post]
func SetUserAuthority(c *gin.Context) {
	var sua request.SetUserAuth
	_ = c.ShouldBindJSON(&sua)
	if UserVerifyErr := utils.Verify(sua, utils.SetUserAuthorityVerify); UserVerifyErr != nil {
		response.FailWithMessage(UserVerifyErr.Error(), c)
		return
	}

	err := service.SetUserAuthority(sua.UUID, sua.AuthorityId)
	if err != nil {
		global.GVA_LOG.Error("修改失败!", zap.Any("err", err))
		response.FailWithMessage("修改失败", c)
		return
	}
	response.OkWithMessage("修改成功", c)
}

// @Tags SysUser
// @Summary 用户修改密码
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body request.ChangePasswordStruct true "用户名, 原密码, 新密码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /user/changePassword [put]
func ChangePassword(c *gin.Context) {
	var user request.ChangePasswordRequest
	_ = c.ShouldBindJSON(&user)
	if err := utils.Verify(user, utils.ChangePasswordVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	u := &model.SysUser{Username: user.Username, Password: user.Password}
	err, _ := service.ChangePassword(u, user.NewPassword)
	if err != nil {
		global.GVA_LOG.Error("修改失败!", zap.Any("err", err))
		response.FailWithMessage("修改失败，原密码与当前账户不符", c)
		return
	}
	response.OkWithMessage("修改成功", c)
}

// @Tags SysUser
// @Summary 设置用户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysUser true "ID, 用户名, 昵称, 头像链接"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
// @Router /user/setUserInfo [put]
func SetUserInfo(c *gin.Context) {
	var user model.SysUser
	_ = c.ShouldBindJSON(&user)
	if err := utils.Verify(user, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err, userInter := service.SetUserInfo(&user)
	if err != nil {
		global.GVA_LOG.Error("设置失败!", zap.Any("err", err))
		response.FailWithMessage("设置失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"userInfo": userInter}, "设置成功", c)
}

// 从 Gin 的 Context 中获取从 jwt 解析出来的用户角色 id
func getUserAuthorityId(c *gin.Context) string {
	claims, exists := c.Get("claims")
	if !exists {
		global.GVA_LOG.Error("从 Gin 的 Context 中获取从 jwt 解析出来的用户 AuthorityId 失败, 请检查路由是否使用 jwt 中间件!")
		return ""
	}
	return claims.(*request.CustomClaims).AuthorityId
}

func getUserID(c *gin.Context) uint {
	claims, exists := c.Get("claims")
	if !exists {
		return 0
	}
	return claims.(*request.CustomClaims).ID
}
