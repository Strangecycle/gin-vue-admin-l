package service

import (
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model"
	"gin-vue-admin-l/model/request"
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
	err = db.Limit(limit).Offset(offset).Preload("DataAuthorityId").Where("parent_id = 0").Find(&authority).Error
	if len(authority) > 0 {
		for k := range authority {
			err = findChildrenAuthority(&authority[k])
		}
	}
	return err, authority, total
}
