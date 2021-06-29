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
		m.POST("getMenu", v1.GetMenu)                   // 获取菜单树
		m.POST("getBaseMenuTree", v1.GetBaseMenuTree)   // 获取用户动态路由
		m.POST("getMenuAuthority", v1.GetMenuAuthority) // 获取指定角色menu
		m.POST("addMenuAuthority", v1.AddMenuAuthority) //	增加menu和角色关联关系
	}
}
