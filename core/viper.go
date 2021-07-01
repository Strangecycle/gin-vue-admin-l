package core

import (
	"flag"
	"fmt"
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

func Viper(path ...string) *viper.Viper {
	config := ""
	if len(path) == 0 {
		flag.StringVar(&config, "c", "", "choose config file.")
		flag.Parse()
		if config == "" { // 优先级: 命令行 > 环境变量 > 默认值
			if configEnv := os.Getenv(utils.ConfigEnv); configEnv == "" {
				config = utils.ConfigFile
				fmt.Printf("您正在使用 config 的默认值,config 的路径为%v\n", utils.ConfigFile)
			} else {
				config = configEnv
				fmt.Printf("您正在使用 GVA_CONFIG 环境变量, config 的路径为%v\n", config)
			}
		} else {
			fmt.Printf("您正在使用命令行的-c参数传递的值,config 的路径为%v\n", config)
		}
	} else {
		config = path[0]
		fmt.Printf("您正在使用 func Viper() 传递的值，config 的路径为 %v\n", config)
	}

	// 初始化 viper
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig() // 读取配置文件
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	// 配置文件内容序列化赋值给全局配置对象
	if err := v.Unmarshal(&global.GVA_CONFIG); err != nil {
		fmt.Println(err)
	}

	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config file changed:", in.Name)
		// 监听配置文件改变序列化配置对象
		if err := v.Unmarshal(&global.GVA_CONFIG); err != nil {
			fmt.Println(err)
		}
	})

	// TODO global.GVA_CONFIG.AutoCode.Root

	return v
}
