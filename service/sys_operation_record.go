package service

import (
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model"
)

func CreateSysOperationRecord(r model.SysOperationRecord) (err error) {
	err = global.GVA_DB.Create(&r).Error
	return err
}
