package request

import "gin-vue-admin-l/model"

type SearchApiParams struct {
	model.SysApi
	PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     string `json:"desc"`     // 描述
}
