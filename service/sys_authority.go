package service

import (
	"errors"
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model"
	"gin-vue-admin-l/model/request"
	"gin-vue-admin-l/model/response"
	"gorm.io/gorm"
	"strconv"
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

func UpdateAuthority(auth model.SysAuthority) (err error, authority model.SysAuthority) {
	err = global.GVA_DB.Where("authority_id = ?", auth.AuthorityId).First(&authority).Updates(&auth).Error
	return err, auth
}

// 流程：
// 1、有用户是该角色或该角色存在子角色时无法删除
// 2、如果存在软删除记录，永久删除它们
// 3、解除角色和菜单的关系，不删除数据库中的菜单关联数据
// 4、删除该角色的所有权限
func DeleteAuthority(auth *model.SysAuthority) (err error) {
	if !errors.Is(global.GVA_DB.Where("authority_id = ?", auth.AuthorityId).First(&model.SysUser{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("此角色有用户正在使用禁止删除")
	}
	if !errors.Is(global.GVA_DB.Where("parent_id = ?", auth.AuthorityId).First(&model.SysAuthority{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("此角色存在子角色不允许删除")
	}

	db := global.GVA_DB.Where("authority_id = ?", auth.AuthorityId).Preload("SysBaseMenus").First(auth)
	// 将软删除的记录永久删除
	err = db.Unscoped().Delete(auth).Error
	if len(auth.SysBaseMenus) > 0 {
		// 如果存在，则删除源模型与参数之间的关系，只会删除引用，不会从数据库中删除这些对象（不会将 sys_base_menus 中的表删除）
		err = global.GVA_DB.Model(auth).Association("SysBaseMenus").Delete(auth.SysBaseMenus)
	} else {
		err = db.Error
	}

	// 删除该角色的权限
	ClearCasbin(0, auth.AuthorityId)
	return err
}

func CopyAuthority(infoCopy response.SysAuthorityCopyResponse) (err error, authority model.SysAuthority) {
	if !errors.Is(global.GVA_DB.Where("authority_id = ?", infoCopy.Authority.AuthorityId).First(&authority).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同角色id"), authority
	}

	infoCopy.Authority.Children = []model.SysAuthority{}
	// 从 authority_menu 视图找出被拷贝角色对应的菜单
	err, menus := GetMenuAuthority(&request.GetAuthorityId{AuthorityId: infoCopy.OldAuthorityId})
	var baseMenu []model.SysBaseMenu
	for _, v := range menus {
		intNum, _ := strconv.Atoi(v.MenuId)
		v.SysBaseMenu.ID = uint(intNum)
		baseMenu = append(baseMenu, v.SysBaseMenu)
	}
	// 将被拷贝角色菜单处理过后赋值给新角色的菜单
	infoCopy.Authority.SysBaseMenus = baseMenu
	err = global.GVA_DB.Create(&infoCopy.Authority).Error

	// 拿到被拷贝角色的所有权限（策略）
	paths := GetPolicyPathByAuthorityId(infoCopy.OldAuthorityId)
	// 给新角色赋予这些权限
	err = UpdateCasbin(infoCopy.Authority.AuthorityId, paths)
	// 出现错误则删除这个角色
	if err != nil {
		_ = DeleteAuthority(&infoCopy.Authority)
	}
	return err, infoCopy.Authority
}
