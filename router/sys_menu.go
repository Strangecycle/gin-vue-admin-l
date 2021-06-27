package router

import (
	v1 "gin-vue-admin-l/api/v1"
	"gin-vue-admin-l/middleware"
	"github.com/gin-gonic/gin"
)

func InitMenuRouter(rg *gin.RouterGroup) {
	m := rg.Group("menu")
	m.Use(middleware.OperationRecord())
	{
		m.POST("getMenu", v1.GetMenu) // 获取菜单树
	}
}
