package core

import (
	"fmt"
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/initialize"
	"go.uber.org/zap"
	"time"
)

type server interface {
	ListenAndServe() error
}

// 启动 http 服务
func RunWindowsServer() {
	// 初始化 redis 服务
	if global.GVA_CONFIG.System.UseMultipoint {
		initialize.Redis()
	}

	r := initialize.Routers()
	// TODO form-generator 开放静态资源

	addr := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	s := initServer(addr, r)

	// 保证文本顺序输出
	time.Sleep(10 * time.Microsecond)
	global.GVA_LOG.Info("Server is running on ", zap.String("address", addr))

	fmt.Printf(`
		当前版本:V2.4.2
		默认自动化文档地址:http://127.0.0.1%s/swagger/index.html
		默认前端文件运行地址:http://127.0.0.1:8080
	`, addr)

	global.GVA_LOG.Error(s.ListenAndServe().Error())
}
