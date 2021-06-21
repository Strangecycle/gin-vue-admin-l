package router

import (
	v1 "gin-vue-admin-l/api/v1"
	"github.com/gin-gonic/gin"
)

func InitInitRouter(pg *gin.RouterGroup) {
	i := pg.Group("init")
	{
		i.POST("checkdb", v1.CheckDB) // 检查是否初始化数据库
		i.POST("initdb", v1.InitDB)   // 初始化数据库配置
	}
}
