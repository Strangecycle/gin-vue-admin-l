package response

import "gin-vue-admin-l/model"

type SysAuthorityResponse struct {
	Authority model.SysAuthority `json:"authority"`
}
