package main

import (
	"gin-vue-admin-l/core"
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/initialize"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

func main() {
	global.GVA_VP = core.Viper()      // 初始化 viper
	global.GVA_LOG = core.Zap()       // 初始化 zap 日志库
	global.GVA_DB = initialize.Gorm() // 初始化数据库连接
	// 初始化定时任务，用于定时清理数据库表数据（如操作记录）
	initialize.Timer()

	if global.GVA_DB != nil {
		db, _ := global.GVA_DB.DB()
		// 程序退出时断开数据库连接
		defer db.Close()
	}

	// 运行服务
	core.RunWindowsServer()
}
