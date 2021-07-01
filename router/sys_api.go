package router

import (
	v1 "gin-vue-admin-l/api/v1"
	"gin-vue-admin-l/middleware"
	"github.com/gin-gonic/gin"
)

func InitApiRouter(rg *gin.RouterGroup) {
	a := rg.Group("api")
	a.Use(middleware.OperationRecord())
	{
		a.POST("getAllApis", v1.GetAllApis) // 获取所有 api
	}
}
