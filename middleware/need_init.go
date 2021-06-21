package middleware

import (
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model/response"
	"github.com/gin-gonic/gin"
)

func NeedInit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 判断是否通过前端配置过数据库信息
		if global.GVA_DB == nil {
			response.OkWithDetailed(gin.H{
				"needInit": true,
			}, "前往初始化数据库", c)
			c.Abort()
		} else {
			c.Next()
		}
	}
}
