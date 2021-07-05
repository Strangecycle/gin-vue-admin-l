package service

import (
	"errors"
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model"
	"gin-vue-admin-l/model/request"
	"gorm.io/gorm"
)

func GetAllApis() (err error, api []model.SysApi) {
	var apis []model.SysApi
	err = global.GVA_DB.Find(&apis).Error
	return err, apis
}

func GetApiListInfoList(info request.SearchApiParams) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&model.SysApi{})
	var apiList []model.SysApi

	// 过滤筛选条件
	if info.Path != "" {
		db = db.Where("path LIKE ?", "%"+info.Path+"%")
	}
	if info.Description != "" {
		db = db.Where("description LIKE ?", "%"+info.Description+"%")
	}
	if info.Method != "" {
		db = db.Where("method LIKE ?", "%"+info.Method+"%")
	}
	if info.ApiGroup != "" {
		db = db.Where("api_group = ?", info.ApiGroup)
	}

	err = db.Count(&total).Error
	if err != nil {
		return err, apiList, total
	}

	db = db.Limit(limit).Offset(offset)
	if info.OrderKey != "" {
		orderStr := ""
		if info.Desc != "" {
			orderStr = info.OrderKey + " desc"
		} else {
			orderStr = info.OrderKey
		}
		err = db.Order(orderStr).Find(&apiList).Error
	} else {
		err = db.Order("api_group").Find(&apiList).Error
	}
	return err, apiList, total
}

func CreateApi(api model.SysApi) (err error) {
	if !errors.Is(global.GVA_DB.Where("path = ? AND method = ?", api.Path, api.Method).First(&model.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同 api")
	}

	err = global.GVA_DB.Create(&api).Error
	return err
}

func GetApiById(id float64) (err error, api model.SysApi) {
	err = global.GVA_DB.Where("id = ?", id).First(&api).Error
	return err, api
}

func UpdateApi(api model.SysApi) (err error) {
	var oldApi model.SysApi
	err = global.GVA_DB.Where("id = ?", api.ID).First(&oldApi).Error
	if err != nil {
		return err
	}

	if oldApi.Path != api.Path || oldApi.Method != api.Method {
		if !errors.Is(global.GVA_DB.Where("path = ? AND method = ?", api.Path, api.Method).First(&model.SysApi{}).Error, gorm.ErrRecordNotFound) {
			return errors.New("存在相同api路径")
		}
	}

	// 更新 casbin 表
	err = UpdateCasbinApi(oldApi.Path, api.Path, oldApi.Method, api.Method)
	if err != nil {
		return err
	}

	// 更新 api 表
	err = global.GVA_DB.Save(&api).Error
	return err
}

func DeleteApi(api model.SysApi) (err error) {
	// 软删除
	err = global.GVA_DB.Delete(&api).Error
	ClearCasbin(1, api.Path, api.Method)
	return err
}

func DeleteApisByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Where("id IN (?)", ids.Ids).Delete(&model.SysApi{}).Error
	// TODO 批量删除 casbin 表中的数据
	return err
}
