package middleware

import (
	"github.com/gin-gonic/gin"
)

// Casbin 权限管理中间件
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// claims, _ := c.Get("claims")
		// waitUse := claims.(*request.CustomClaims)
		// // 获取请求的 URI
		// uri := c.Request.URL.RequestURI()
		// // 获取请求方法
		// act := c.Request.Method
		// // 获取用户角色
		// sub := waitUse.AuthorityId
		// service.Casbin()

		c.Next()
	}
}
