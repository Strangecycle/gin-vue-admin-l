package router

import (
	v1 "gin-vue-admin-l/api/v1"
	"github.com/gin-gonic/gin"
)

func InitJwtRouter(pg *gin.RouterGroup) {
	// TODO use record middleware
	j := pg.Group("jwt")
	{
		j.POST("jsonInBlacklist", v1.JsonInBlacklist) // jwt 加入黑名单(退出登录)
	}
}
