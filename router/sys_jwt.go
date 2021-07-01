package router

import (
	v1 "gin-vue-admin-l/api/v1"
	"gin-vue-admin-l/middleware"
	"github.com/gin-gonic/gin"
)

func InitJwtRouter(pg *gin.RouterGroup) {
	j := pg.Group("jwt")
	j.Use(middleware.OperationRecord())
	{
		j.POST("jsonInBlacklist", v1.JsonInBlacklist) // jwt 加入黑名单(退出登录)
	}
}
