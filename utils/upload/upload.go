package upload

import (
	"gin-vue-admin-l/global"
	"mime/multipart"
)

type OSS interface {
	UploadFile(file *multipart.FileHeader) (string, string, error)
	DeleteFile(key string) error
}

func NewOss() OSS {
	switch global.GVA_CONFIG.System.OssType {
	case "local":
		return &Local{}
	// TODO 文件上传阿里云 OSS、腾讯云 OSS、七牛云 OSS
	default:
		return &Local{}
	}
}
