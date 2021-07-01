package service

import (
	"database/sql"
	"fmt"
	"gin-vue-admin-l/config"
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model"
	"gin-vue-admin-l/model/request"
	"gin-vue-admin-l/source"
	"gin-vue-admin-l/utils"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(conf request.InitDB) error {
	baseConf := config.Mysql{
		Path:     "",
		Dbname:   "",
		Username: "",
		Password: "",
		Config:   "charset=utf8mb4&parseTime=True&loc=Local",
	}

	// 默认配置
	if conf.Host == "" {
		conf.Host = "127.0.0.1"
	}
	if conf.Port == "" {
		conf.Port = "3306"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", conf.UserName, conf.Password, conf.Host, conf.Port)
	// 初始化创建数据库
	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", conf.DBName)
	err := createTable(dsn, global.GVA_CONFIG.System.DbType, createSql)
	if err != nil {
		return err
	}

	mysqlConf := config.Mysql{
		Path:     fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		Dbname:   conf.DBName,
		Username: conf.UserName,
		Password: conf.Password,
		Config:   "charset=utf8mb4&parseTime=True&loc=Local",
	}

	// 向配置文件中写入数据库配置
	if err := writeConfig(global.GVA_VP, mysqlConf); err != nil {
		return err
	}

	// 由于初始化 viper 时开启了配置文件监听，配置文件改变时会将新的配置赋值到全局配置对象
	m := global.GVA_CONFIG.Mysql
	if m.Dbname == "" {
		return nil
	}

	linkDns := m.Username + ":" + m.Password + "@tcp(" + m.Path + ")/" + m.Dbname + "?" + m.Config
	mConf := mysql.Config{
		DSN:                       linkDns, // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		DisableDatetimePrecision:  true,    // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,    // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,    // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	db, err := gorm.Open(mysql.New(mConf), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		// 连接 mysql 失败，写入默认配置
		_ = writeConfig(global.GVA_VP, baseConf)
		return nil
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(m.MaxIdleConns)
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)
	global.GVA_DB = db

	// 自动创建表
	err = global.GVA_DB.AutoMigrate(
		model.SysUser{},
		model.SysAuthority{},
		model.SysBaseMenu{},
		model.SysBaseMenuParameter{},
		model.JwtBlacklist{},
		model.SysOperationRecord{},
		model.SysApi{},
	)
	if err != nil {
		_ = writeConfig(global.GVA_VP, baseConf)
		return err
	}

	// initDB() 向数据表中填入初始数据
	initDB(
		source.Admin,
		source.Authority,
		source.DataAuthorities,
		source.BaseMenu,
		source.AuthoritiesMenus,
		source.AuthorityMenu, // 视图
		source.Casbin,        // casbin 策略表
		source.Api,
	)

	// TODO global.GVA_CONFIG.AutoCode.Root

	return nil
}

// 创建数据库
func createTable(dsn string, driver string, createSql string) error {
	// 官方标准库数据库连接
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// Ping 验证与数据库的连接是否仍然有效，必要时建立连接。
	if err := db.Ping(); err != nil {
		return err
	}

	// 执行 sql
	_, err = db.Exec(createSql)
	return err
}

// 向配置文件中写入配置
func writeConfig(viper *viper.Viper, mysqlConf config.Mysql) error {
	global.GVA_CONFIG.Mysql = mysqlConf
	confMap := utils.StuctToMap(global.GVA_CONFIG)
	for k, v := range confMap {
		viper.Set(k, v)
	}
	return viper.WriteConfig()
}

func initDB(initDBFunc ...model.InitDBFunc) (err error) {
	for _, v := range initDBFunc {
		if err := v.Init(); err != nil {
			return err
		}
	}
	return nil
}
