package initialize

import (
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/middleware"
	"gin-vue-admin-l/router"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 初始化总路由
func Routers() *gin.Engine {
	r := gin.Default()
	// 为用户头像和文件提供静态地址
	r.StaticFS(global.GVA_CONFIG.Local.Path, http.Dir(global.GVA_CONFIG.Local.Path))
	global.GVA_LOG.Info("use middleware logger")
	// 跨域
	// r.Use(middleware.Cors()) // 如需跨域可以打开
	global.GVA_LOG.Info("use middleware cors")

	// TODO swagger 路由
	global.GVA_LOG.Info("register swagger handler")

	// 公共路由（无需鉴权）
	publicGroup := r.Group("")
	{
		router.InitInitRouter(publicGroup) // 自动初始化相关
		router.InitBaseRouter(publicGroup) // 注册基础功能路由
	}

	// 需要鉴权的路由
	privateGroup := r.Group("")
	privateGroup.Use(middleware.JWTAuth(), middleware.CasbinHandler())
	{
		router.InitJwtRouter(privateGroup) // jwt相关路由
	}
	global.GVA_LOG.Info("router register success")
	return r
}
