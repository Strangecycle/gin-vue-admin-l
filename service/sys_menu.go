package service

import (
	"errors"
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model"
	"gin-vue-admin-l/model/request"
	"gorm.io/gorm"
	"strconv"
)

// 获取路由总树 map
func getMenuTreeMap(authId string) (err error, treeMap map[string][]model.SysMenu) {
	var allMenus []model.SysMenu
	treeMap = make(map[string][]model.SysMenu)
	err = global.GVA_DB.Where("authority_id = ?", authId).Order("sort").Preload("Parameters").Find(&allMenus).Error
	for _, v := range allMenus {
		// 每一条菜单都有一个 id 和 parentId，当 parentId 为 0 时，表示此路由是一级路由
		/*{
			'0': [{ id: '1', path: 'aboutus' }], // 一级路由
			'9': [{ id: '22', path: 'foo' }, { id: '23', path: 'bar' }] // id 为 9 的路由的子路由
		}*/
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}
	return err, treeMap
}

// 获取子菜单
func getChildrenList(menu *model.SysMenu, treeMap map[string][]model.SysMenu) (err error) {
	// treeMap 已经将有子路由的数据映射出来了，所以直接根据当前 menuId 取 map 中的值赋给当前路由的 children
	// menuId 未命中 map 中的值，则表示当前无子路由，返回 nil
	menu.Children = treeMap[menu.MenuId]
	for i := 0; i < len(menu.Children); i++ {
		err = getChildrenList(&menu.Children[i], treeMap)
	}
	return err
}

// 获取菜单树
func GetMenuTree(authId string) (err error, menus []model.SysMenu) {
	err, menuTree := getMenuTreeMap(authId)
	// 拿到所有一级路由
	menus = menuTree["0"]
	// 获取一级路由的子路由
	for i := 0; i < len(menus); i++ {
		err = getChildrenList(&menus[i], menuTree)
	}
	return err, menus
}

// 获取路由总树 map
func getBaseMenuTreeMap() (err error, treeMap map[string][]model.SysBaseMenu) {
	var allMenus []model.SysBaseMenu
	treeMap = make(map[string][]model.SysBaseMenu)
	err = global.GVA_DB.Order("sort").Preload("Parameters").Find(&allMenus).Error
	for _, v := range allMenus {
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}
	return err, treeMap
}

// 获取菜单的子菜单
func getBaseChildrenList(menu *model.SysBaseMenu, treeMap map[string][]model.SysBaseMenu) (err error) {
	menu.Children = treeMap[strconv.Itoa(int(menu.ID))]
	for i := 0; i < len(menu.Children); i++ {
		err = getBaseChildrenList(&menu.Children[i], treeMap)
	}
	return err
}

// 获取基础路由树
func GetBaseMenuTree() (err error, menus []model.SysBaseMenu) {
	err, treeMap := getBaseMenuTreeMap()
	menus = treeMap["0"]
	for i := 0; i < len(menus); i++ {
		err = getBaseChildrenList(&menus[i], treeMap)
	}
	return err, menus
}

// 获取路由列表
func GetMenuList() (err error, list interface{}, total int64) {
	var menuList []model.SysBaseMenu
	err, treeMap := getBaseMenuTreeMap()
	menuList = treeMap["0"]
	for i := 0; i < len(menuList); i++ {
		err = getBaseChildrenList(&menuList[i], treeMap)
	}
	return err, menuList, total
}

func GetMenuAuthority(param *request.GetAuthorityId) (err error, menus []model.SysMenu) {
	err = global.GVA_DB.Where("authority_id = ?", param.AuthorityId).Order("sort").Find(&menus).Error
	return err, menus
}

func AddMenuAuthority(menus []model.SysBaseMenu, authId string) (err error) {
	var auth model.SysAuthority
	auth.AuthorityId = authId
	auth.SysBaseMenus = menus
	err = SetMenuAuthority(&auth)
	return err
}

func AddBaseMenu(menu model.SysBaseMenu) (err error) {
	if !errors.Is(global.GVA_DB.Where("name = ?", menu.Name).First(&model.SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复name，请修改name")
	}
	err = global.GVA_DB.Create(&menu).Error
	return err
}

func GetBaseMenuById(id float64) (err error, menu model.SysBaseMenu) {
	err = global.GVA_DB.Preload("Parameters").Where("id = ?", id).First(&menu).Error
	return err, menu
}

func UpdateBaseMenu(menu model.SysBaseMenu) (err error) {
	var oldMenu model.SysBaseMenu
	updateMap := make(map[string]interface{})
	updateMap["keep_alive"] = menu.KeepAlive
	updateMap["default_menu"] = menu.DefaultMenu
	updateMap["parent_id"] = menu.ParentId
	updateMap["path"] = menu.Path
	updateMap["name"] = menu.Name
	updateMap["hidden"] = menu.Hidden
	updateMap["component"] = menu.Component
	updateMap["title"] = menu.Title
	updateMap["icon"] = menu.Icon
	updateMap["sort"] = menu.Sort

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", menu.ID).Find(&oldMenu)
		if oldMenu.Name != menu.Name {
			if !errors.Is(tx.Where("id <> ? AND name = ?", menu.ID, menu.ID).First(&model.SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
				global.GVA_LOG.Debug("存在相同 name 修改失败")
				return errors.New("存在相同 name 修改失败")
			}
		}

		// 1、永久删除这个路由的参数
		txErr := tx.Unscoped().Where("sys_base_menu_id = ?", menu.ID).Delete(&model.SysBaseMenuParameter{}).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}

		// 2、存入新的路由参数
		if len(menu.Parameters) > 0 {
			for k, _ := range menu.Parameters {
				menu.Parameters[k].SysBaseMenuID = menu.ID
			}
			// 路由参数存入数据库
			txErr := tx.Create(&menu.Parameters).Error
			if txErr != nil {
				global.GVA_LOG.Debug(txErr.Error())
				return txErr
			}
		}

		// 3、更新菜单
		// 这里使用 db 是因为上面 db 已经绑定了 Menu 模型
		txErr = db.Updates(updateMap).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}

		return nil
	})

	return err
}

func DeleteBaseMenu(id float64) (err error) {
	// 找出该菜单的子菜单
	err = global.GVA_DB.Preload("Parameters").Where("parent_id = ?", id).First(&model.SysBaseMenu{}).Error
	if err != nil {
		var menu model.SysBaseMenu
		// 先软删除这个菜单
		db := global.GVA_DB.Preload("SysAuthoritys").Where("id = ?", id).First(&menu).Delete(&menu)
		// 然后删除这个菜单的参数
		err = global.GVA_DB.Where("sys_base_menu_id = ?", id).Delete(&model.SysBaseMenuParameter{}).Error
		if len(menu.SysAuthoritys) > 0 {
			// 解除此菜单关联的所有角色的关联，否则它们之间存在引用无法删除
			err = global.GVA_DB.Model(&menu).Association("SysAuthoritys").Delete(&menu.SysAuthoritys)
			return err
		}
		return db.Error
	}
	return errors.New("此菜单存在子菜单不可删除")
}
