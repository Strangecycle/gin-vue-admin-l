package config

// 整合服务端工具的配置对象
type Server struct {
	System  System  `mapstructure:"system" json:"system" yaml:"system"`
	Local   Local   `mapstructure:"local" json:"local" yaml:"local"`
	Zap     Zap     `mapstructure:"zap" json:"zap" yaml:"zap"`
	Mysql   Mysql   `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Timer   Timer   `mapstructure:"timer" json:"timer" yaml:"timer"`
	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	JWT     JWT     `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Redis   Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
	Casbin  Casbin  `mapstructure:"casbin" json:"casbin" yaml:"casbin"`
	// 示例模块
	Excel Excel `mapstructure:"excel" json:"excel" yaml:"excel"`
}
