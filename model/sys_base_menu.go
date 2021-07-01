package model

import "gin-vue-admin-l/global"

// 对应前端 Vue 路由表
type SysBaseMenu struct {
	global.GVA_MODEL
	MenuLevel     uint                              `json:"-"`
	ParentId      string                            `json:"parentId" gorm:"comment:父菜单ID"`     // 父菜单ID
	Path          string                            `json:"path" gorm:"父菜单ID"`                 // 父菜单ID
	Name          string                            `json:"name" gorm:"comment:路由name"`        // 路由name
	Hidden        bool                              `json:"hidden" gorm:"comment:是否在列表隐藏"`     // 是否在列表隐藏
	Component     string                            `json:"component" gorm:"comment:对应前端文件路径"` // 对应前端文件路径
	Sort          int                               `json:"sort" gorm:"comment:排序标记"`          // 排序标记
	Meta          `json:"meta" gorm:"comment:附加属性"` // 附加属性
	SysAuthoritys []SysAuthority                    `json:"sysAuthoritys" gorm:"many2many:sys_authority_menus"` // 权限表与菜单表多对多关系，创建 sys_authority_menus 中间表
	Children      []SysBaseMenu                     `json:"children" gorm:"-"`
	Parameters    []SysBaseMenuParameter            `json:"parameters"`
}

// 路由元数据
type Meta struct {
	KeepAlive   bool   `json:"keepAlive" gorm:"comment:是否缓存"`           // 是否缓存
	DefaultMenu bool   `json:"defaultMenu" gorm:"comment:是否是基础路由（开发中）"` // 是否是基础路由（开发中）
	Title       string `json:"title" gorm:"comment:菜单名"`                // 菜单名
	Icon        string `json:"icon" gorm:"comment:菜单图标"`                // 菜单图标
	CloseTab    bool   `json:"closeTab" gorm:"comment:自动关闭tab"`         // 自动关闭tab
}

// 路由参数
type SysBaseMenuParameter struct {
	global.GVA_MODEL
	SysBaseMenuID uint
	Type          string `json:"type" gorm:"comment:地址栏携带参数为params还是query"` // 地址栏携带参数为params还是query
	Key           string `json:"key" gorm:"comment:地址栏携带参数的key"`            // 地址栏携带参数的key
	Value         string `json:"value" gorm:"comment:地址栏携带参数的值"`            // 地址栏携带参数的值
}
