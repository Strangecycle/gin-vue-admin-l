package router

import (
	v1 "gin-vue-admin-l/api/v1"
	"gin-vue-admin-l/middleware"
	"github.com/gin-gonic/gin"
)

func InitBaseRouter(pg *gin.RouterGroup) {
	b := pg.Group("base")
	b.Use(middleware.NeedInit())
	{
		b.POST("captcha", v1.Captcha)
		b.POST("login", v1.Login)
	}
}
