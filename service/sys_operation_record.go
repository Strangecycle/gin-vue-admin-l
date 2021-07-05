package service

import (
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model"
	"gin-vue-admin-l/model/request"
)

func CreateSysOperationRecord(r model.SysOperationRecord) (err error) {
	err = global.GVA_DB.Create(&r).Error
	return err
}

func GetSysOperationRecordInfoList(info request.SysOperationRecordSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	var records []model.SysOperationRecord
	db := global.GVA_DB.Model(&model.SysOperationRecord{})

	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Method != "" {
		db = db.Where("method = ?", info.Method)
	}
	if info.Path != "" {
		db = db.Where("path LIKE ？", "%"+info.Path+"%")
	}
	if info.Status != 0 {
		db = db.Where("status = ?", info.Status)
	}

	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Preload("User").Order("id desc").Find(&records).Error
	return err, records, total
}

func DeleteSysOperationRecord(record model.SysOperationRecord) (err error) {
	err = global.GVA_DB.Delete(&record).Error
	return err
}

func DeleteSysOperationRecordByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Where("id IN (?)", ids.Ids).Delete(&model.SysOperationRecord{}).Error
	return err
}
