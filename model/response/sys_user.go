package response

import "gin-vue-admin-l/model"

type LoginResponse struct {
	User      model.SysUser `json:"user"`
	Token     string        `json:"token"`
	ExpiresAt int64         `json:"expiresAt"` // token 过期时间
}
