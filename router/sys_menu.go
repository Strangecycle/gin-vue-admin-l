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
		m.POST("getMenuAuthority", v1.GetMenuAuthority) // 获取指定角色 menu
		m.POST("addMenuAuthority", v1.AddMenuAuthority) // 增加 menu 和角色关联关系
		m.POST("getMenuList", v1.GetMenuList)           // 获取菜单列表
		m.POST("addBaseMenu", v1.AddBaseMenu)           // 新增菜单
		m.POST("getBaseMenuById", v1.GetBaseMenuById)   // 根据 id 获取菜单
		m.POST("updateBaseMenu", v1.UpdateBaseMenu)     // 更新菜单
		m.POST("deleteBaseMenu", v1.DeleteBaseMenu)     // 删除菜单
	}
}
