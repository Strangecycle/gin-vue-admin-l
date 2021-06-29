package request

import "gin-vue-admin-l/model"

type AddMenuAuthorityInfo struct {
	AuthorityId string              `json:"authorityId"`
	Menus       []model.SysBaseMenu `json:"menus"`
}
