package service

import (
	"errors"
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model"
	"gin-vue-admin-l/model/request"
	"gorm.io/gorm"
)

// 查询子角色
func findChildrenAuthority(auth *model.SysAuthority) (err error) {
	err = global.GVA_DB.Preload("DataAuthorityId").Where("parent_id = ?", auth.AuthorityId).Find(&auth.Children).Error
	if len(auth.Children) > 0 {
		for k := range auth.Children {
			err = findChildrenAuthority(&auth.Children[k])
		}
	}
	return err
}

func GetAuthorityInfoList(p request.PageInfo) (err error, list interface{}, total int64) {
	limit := p.PageSize
	offset := p.PageSize * (p.Page - 1)
	db := global.GVA_DB.Model(&model.SysAuthority{})
	var authority []model.SysAuthority
	// 凡是 many2many 关系的表，Preload 的时候会先拿到当前表的主键到中间表去找匹配 Preload 的表的主键，拿到对应主键后再到 Preload 绑定的表中去查找
	// 1、select `authority_id` from `sys_authorities`
	// 2、select `data_authority_id_authority_id` from `sys_data_authority_id` where `sys_authority_authority_id` IN (authority_id...)
	// 3、model.SysAuthority.DataAuthorityId = select * from `sys_authorities` where `authority_id` IN (data_authority_id_authority_id...)
	err = db.Limit(limit).Offset(offset).Preload("DataAuthorityId").Where("parent_id = 0").Find(&authority).Error
	if len(authority) > 0 {
		for k := range authority {
			err = findChildrenAuthority(&authority[k])
		}
	}
	return err, authority, total
}

func CreateAuthority(auth model.SysAuthority) (err error, authInter model.SysAuthority) {
	var authority model.SysAuthority
	err = global.GVA_DB.Where("authority_id = ?", auth.AuthorityId).First(&authority).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("存在相同角色id"), auth
	}
	err = global.GVA_DB.Create(&auth).Error
	return err, auth
}

// 菜单与角色绑定
// 流程：
// 1、从 authority 表中找到这个角色
// 2、通过 authority_menu 中间表找到角色对应的菜单
// 3、利用 gorm 的关联模式将用用户传进来的菜单替换角色的旧菜单
func SetMenuAuthority(auth *model.SysAuthority) error {
	var s model.SysAuthority
	// 这里是到 sys_authority_menus 这个中间表根据 authority_id 去查找 base_menu_id，然后 Preload 再拿 base_menu_id 去预加载（根据 id 到 sys_base_menus 找）到 SysAuthority.SysBaseMenus
	// 1、select `authority_id` from `sys_authorities`
	// 2、select `id` from `sys_authority_menus` where `sys_authority_authority_id` IN (authority_id...)
	// 3、s.SysBaseMenus = select * from sys_base_menus where id IN (id...)
	global.GVA_DB.Preload("SysBaseMenus").First(&s, "authority_id = ?", auth.AuthorityId)
	// 这行调用会先找出数据库中 s 与传进来的 auth 相匹配的记录，然后将传进来的 auth.SysBaseMenus 替换掉数据库中的 s.SysBaseMenus
	// 以此达到更新角色菜单的目的
	err := global.GVA_DB.Model(&s).Association("SysBaseMenus").Replace(&auth.SysBaseMenus)
	return err
}

func SetDataAuthority(auth model.SysAuthority) error {
	var s model.SysAuthority
	// global.GVA_DB.Preload("DataAuthorityId").First(&s, "authority_id = ?", auth.AuthorityId)
	global.GVA_DB.Where("authority_id = ?", auth.AuthorityId).Preload("DataAuthorityId").First(&s)
	err := global.GVA_DB.Model(&s).Association("DataAuthorityId").Replace(&auth.DataAuthorityId)
	return err
}
