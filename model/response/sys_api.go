package response

import "gin-vue-admin-l/model"

type SysAPIListResponse struct {
	Apis []model.SysApi `json:"apis"`
}
