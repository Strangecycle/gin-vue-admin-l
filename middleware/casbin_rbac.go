package middleware

import (
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model/request"
	"gin-vue-admin-l/model/response"
	"gin-vue-admin-l/service"
	"github.com/gin-gonic/gin"
)

// Casbin 权限管理中间件
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		waitUse := claims.(*request.CustomClaims)
		// 用户角色
		sub := waitUse.AuthorityId
		// 请求的 uri
		obj := c.Request.URL.RequestURI()
		// 行为
		act := c.Request.Method

		// casbin 模型
		e := service.Casbin()

		// 判断策略中是否存在
		success, _ := e.Enforce(sub, obj, act)
		if global.GVA_CONFIG.System.Env != "develop" || !success {
			response.FailWithDetailed(gin.H{}, "权限不足", c)
			c.Abort()
			return
		}
		c.Next()
	}
}
