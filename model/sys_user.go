package model

import (
	"gin-vue-admin-l/global"
	uuid "github.com/satori/go.uuid"
)

type SysUser struct {
	global.GVA_MODEL
	UUID      uuid.UUID `json:"uuid" gorm:"comment:用户UUID"`                                                    // 用户UUID
	Username  string    `json:"username" gorm:"comment:用户登录名"`                                                 // 用户登录名
	Password  string    `json:"-"  gorm:"comment:用户登录密码"`                                                      // 用户登录密码
	NickName  string    `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`                                     // 用户昵称
	HeaderImg string    `json:"headerImg" gorm:"default:http://qmplusimg.henrongyi.top/head.png;comment:用户头像"` // 用户头像
	// 关联权限表
	Authority   SysAuthority `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"` // 用户角色
	AuthorityId string       `json:"authorityId" gorm:"default:8888;comment:用户角色ID"`                              // 用户角色ID，默认 8888 超管
}
