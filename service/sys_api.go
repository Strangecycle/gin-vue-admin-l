package service

import (
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model"
)

func GetAllApis() (err error, api []model.SysApi) {
	var apis []model.SysApi
	err = global.GVA_DB.Find(&apis).Error
	return err, apis
}
