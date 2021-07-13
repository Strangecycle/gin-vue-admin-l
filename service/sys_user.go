package service

import (
	"errors"
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model"
	"gin-vue-admin-l/model/request"
	"gin-vue-admin-l/utils"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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

func GetUserInfoList(p request.PageInfo) (err error, list interface{}, total int64) {
	limit := p.PageSize
	offset := p.PageSize * (p.Page - 1) // 10 * (1 - 1) -> 10 * (2 - 1)
	db := global.GVA_DB.Model(&model.SysUser{})
	var userList []model.SysUser
	err = db.Count(&total).Error
	// Offset 指定返回记录之前要跳过的记录数，如查询第二页时会查出二十条，offset 为 10，则跳过前十条(前一页)
	err = db.Limit(limit).Offset(offset).Preload("Authority").Find(&userList).Error
	return err, userList, total
}

func Register(u model.SysUser) (err error, userInter model.SysUser) {
	var user model.SysUser
	if !errors.Is(global.GVA_DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) {
		return errors.New("用户名已被注册"), userInter
	}

	u.UUID = uuid.NewV4()
	u.Password = utils.MD5V([]byte(u.Password))
	err = global.GVA_DB.Create(&u).Error

	return err, u
}

func DeleteUser(id float64) (err error) {
	var user model.SysUser
	err = global.GVA_DB.Where("id = ?", id).Delete(&user).Error
	return err
}

func SetUserAuthority(uuid uuid.UUID, authId string) (err error) {
	var user model.SysUser
	err = global.GVA_DB.Where("uuid = ?", uuid).First(&user).Update("authority_id", authId).Error
	return err
}

func ChangePassword(u *model.SysUser, newPwd string) (err error, userInter *model.SysUser) {
	var user model.SysUser
	u.Password = utils.MD5V([]byte(u.Password))
	newPwdMD5 := utils.MD5V([]byte(newPwd))
	err = global.GVA_DB.Where("username = ?", u.Username).First(&user).Update("password", newPwdMD5).Error
	return err, u
}

func SetUserInfo(user *model.SysUser) (err error, userInter model.SysUser) {
	err = global.GVA_DB.Updates(&user).Error
	return err, *user
}
