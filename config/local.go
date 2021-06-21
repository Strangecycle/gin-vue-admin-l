package config

type Local struct {
	Path string `mapstructure:"path" json:"path" yaml:"path"` // 静态资源路径
}
