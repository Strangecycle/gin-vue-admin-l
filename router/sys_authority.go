package router

import (
	v1 "gin-vue-admin-l/api/v1"
	"gin-vue-admin-l/middleware"
	"github.com/gin-gonic/gin"
)

func InitAuthorityRouter(rg *gin.RouterGroup) {
	a := rg.Group("authority")
	a.Use(middleware.OperationRecord())
	{
		a.POST("getAuthorityList", v1.GetAuthorityList) // 获取角色列表
		a.POST("createAuthority", v1.CreateAuthority)   // 新建角色
		a.POST("setDataAuthority", v1.SetDataAuthority) // 设置角色资源权限
	}
}
