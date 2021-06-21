package service

import (
	"errors"
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model"
	"gin-vue-admin-l/utils"
	uuid "github.com/satori/go.uuid"
)

func Login(u *model.SysUser) (err error, userInter *model.SysUser) {
	var user model.SysUser
	// 密码加密与数据库对比
	u.Password = utils.MD5V([]byte(u.Password))
	err = global.GVA_DB.Where("username = ? AND password = ?", u.Username, u.Password).Preload("Authority").First(&user).Error
	return err, &user
}

func FindUserByUuid(uid uuid.UUID) (err error, user *model.SysUser) {
	if err = global.GVA_DB.Where("uuid = ?", uid).First(&user).Error; err != nil {
		return errors.New("用户不存在"), user
	}
	return nil, user
}
