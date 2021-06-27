package service

import (
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model"
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
