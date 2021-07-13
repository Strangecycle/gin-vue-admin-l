package router

import (
	v1 "gin-vue-admin-l/api/v1"
	"gin-vue-admin-l/middleware"
	"github.com/gin-gonic/gin"
)

func InitSystemRouter(rg *gin.RouterGroup) {
	s := rg.Group("system")
	s.Use(middleware.OperationRecord())
	{
		s.POST("getServerInfo", v1.GetServerInfo) // 获取服务器信息
	}
}
