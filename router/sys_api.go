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
		a.POST("getAllApis", v1.GetAllApis)             // 获取所有 api
		a.POST("getApiList", v1.GetApiList)             // 获取 api 列表
		a.POST("createApi", v1.CreateApi)               // 创建 api
		a.POST("getApiById", v1.GetApiById)             // 获取单条 api 消息
		a.POST("updateApi", v1.UpdateApi)               // 编辑 api
		a.POST("deleteApi", v1.DeleteApi)               // 删除 api
		a.DELETE("deleteApisByIds", v1.DeleteApisByIds) // 删除选中api
	}
}
