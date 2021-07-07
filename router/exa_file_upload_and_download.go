package router

import (
	v1 "gin-vue-admin-l/api/v1"
	"github.com/gin-gonic/gin"
)

func InitFileUploadAndDownloadRouter(rg *gin.RouterGroup) {
	f := rg.Group("fileUploadAndDownload")
	{
		f.POST("/upload", v1.UploadFile)      // 上传文件
		f.POST("getFileList", v1.GetFileList) // 获取上传文件列表
		f.POST("deleteFile", v1.DeleteFile)   // 删除上传的文件
	}
}
