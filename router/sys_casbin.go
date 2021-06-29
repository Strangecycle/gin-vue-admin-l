package router

import (
	v1 "gin-vue-admin-l/api/v1"
	"gin-vue-admin-l/middleware"
	"github.com/gin-gonic/gin"
)

func InitCasbinRouter(rg *gin.RouterGroup) {
	c := rg.Group("casbin")
	c.Use(middleware.OperationRecord())
	{
		c.POST("getPolicyPathByAuthorityId", v1.GetPolicyPathByAuthorityId) // 根据 authId 获取权限
		c.POST("updateCasbin", v1.UpdateCasbin)                             // 更新权限
	}
}
