package model

import "time"

// 用户权限表
type SysAuthority struct {
	CreatedAt     time.Time  // 创建时间
	UpdatedAt     time.Time  // 更新时间
	DeletedAt     *time.Time `sql:"index"`
	AuthorityId   string     `json:"authorityId" gorm:"not null;unique;primary_key;comment:角色ID;size:90"` // 角色 ID
	AuthorityName string     `json:"authorityName" gorm:"comment:角色名"`                                    // 角色名
	ParentId      string     `json:"parentId" gorm:"comment:父角色ID"`                                       // 父角色ID
	// 多对多关系，例如：一篇文章属于某个标签（文章表），某个标签下有若干篇文章（标签表），形成多对多关系，相当于形成了一张中间表 tags_articles
	// 这里以 888（超管）、8881（超管子角色）、9528（普通管理员角色）举例
	// 888 包含了 888、8881、9528 三个角色的所有权限
	// 而 9528 只包含了 8881、9528 两个角色的权限
	// sys_data_authority_id 中间表就是为了记录角色之前的包含和所属关系
	DataAuthorityId []SysAuthority `json:"dataAuthorityId" gorm:"many2many:sys_data_authority_id"`
	Children        []SysAuthority `json:"children" gorm:"-"`                                   // - 表示忽略该字段，- 无读写权限，一般用于以树结构返回前端
	SysBaseMenus    []SysBaseMenu  `json:"menus" gorm:"many2many:sys_authority_menus;"`         // 权限表与菜单表多对多关系，创建 sys_authority_menus 中间表
	DefaultRouter   string         `json:"defaultRouter" gorm:"comment:默认菜单;default:dashboard"` // 默认菜单(默认打开的首页)
}
