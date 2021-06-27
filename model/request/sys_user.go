package request

import uuid "github.com/satori/go.uuid"

type Login struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Captcha   string `json:"captcha"`
	CaptchaId string `json:"captchaId"`
}

type Register struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	NickName    string `json:"nickName" gorm:"default:'QMPlusUser'"`
	HeaderImg   string `json:"headerImg" gorm:"default:'http://www.henrongyi.top/avatar/lufu.jpg'"`
	AuthorityId string `json:"authorityId" gorm:"default:'888'"`
}

type SetUserAuth struct {
	UUID        uuid.UUID `json:"uuid"`        // 用户UUID
	AuthorityId string    `json:"authorityId"` // 角色ID
}

type ChangePasswordRequest struct {
	Username    string `json:"username"`    // 用户名
	Password    string `json:"password"`    // 密码
	NewPassword string `json:"newPassword"` // 新密码
}
