package service

import (
	"errors"
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model"
	"gin-vue-admin-l/model/request"
	"gin-vue-admin-l/utils/upload"
	"mime/multipart"
	"strings"
)

func Upload(file model.ExaFileUploadAndDownload) (err error) {
	// 将文件描述对象存入数据库
	return global.GVA_DB.Create(&file).Error
}

func UploadFile(header *multipart.FileHeader, noSave string) (err error, file model.ExaFileUploadAndDownload) {
	oss := upload.NewOss()
	filePath, key, uploadErr := oss.UploadFile(header)
	if uploadErr != nil {
		panic(err)
	}

	if noSave == "0" {
		s := strings.Split(header.Filename, ".")
		file = model.ExaFileUploadAndDownload{
			Url:  filePath,
			Name: header.Filename,
			Tag:  s[len(s)-1],
			Key:  key,
		}
		return Upload(file), file
	}
	return
}

func GetFileRecordInfoList(info request.PageInfo) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	var fileList []model.ExaFileUploadAndDownload
	db := global.GVA_DB
	err = db.Find(&fileList).Count(&total).Error
	err = db.Limit(limit).Offset(offset).Order("created_at desc").Find(&fileList).Error
	return err, fileList, total
}

func FindFile(id uint) (error, model.ExaFileUploadAndDownload) {
	var file model.ExaFileUploadAndDownload
	err := global.GVA_DB.Where("id = ?", id).Find(&file).Error
	return err, file
}

func DeleteFile(file model.ExaFileUploadAndDownload) (err error) {
	var fileFromDb model.ExaFileUploadAndDownload
	err, fileFromDb = FindFile(file.ID)

	oss := upload.NewOss()
	err = oss.DeleteFile(fileFromDb.Key)
	if err != nil {
		return errors.New("删除文件失败")
	}

	// 软删除
	// err = global.GVA_DB.Delete(&file).Error
	// 永久删除
	err = global.GVA_DB.Where("id = ?", file.ID).Unscoped().Delete(&file).Error
	return err
}
