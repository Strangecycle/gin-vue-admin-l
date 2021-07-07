package upload

import (
	"errors"
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/utils"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

type Local struct{}

func (*Local) UploadFile(file *multipart.FileHeader) (string, string, error) {
	// 文件后缀
	ext := path.Ext(file.Filename)
	// 读取文件名并加密
	name := strings.TrimSuffix(file.Filename, ext)
	name = utils.MD5V([]byte(name))
	// 拼接新文件名
	filename := name + "_" + time.Now().Format("20060102150405") + ext
	mkDirErr := os.MkdirAll(global.GVA_CONFIG.Local.Path, os.ModePerm)
	if mkDirErr != nil {
		global.GVA_LOG.Error("function os.MkdirAll() Filed", zap.Any("err", mkDirErr.Error()))
		return "", "", errors.New("function os.MkdirAll() Filed, err:" + mkDirErr.Error())
	}

	// 拼接路径和文件名 -> /uploads/file/xxxx.jpg
	p := global.GVA_CONFIG.Local.Path + "/" + filename

	// 读取上传的文件
	f, openErr := file.Open()
	if openErr != nil {
		global.GVA_LOG.Error("function file.Open() Filed", zap.Any("err", openErr.Error()))
		return "", "", errors.New("function file.Open() Filed, err:" + openErr.Error())
	}
	defer f.Close()

	// 创建新文件
	out, createErr := os.Create(p)
	if createErr != nil {
		global.GVA_LOG.Error("function os.Create() Filed", zap.Any("err", createErr.Error()))
		return "", "", errors.New("function os.Create() Filed, err:" + createErr.Error())
	}
	defer out.Close()

	// 传输（拷贝）文件，将上传的文件内容拷贝到新建的文件中
	_, copyErr := io.Copy(out, f)
	if copyErr != nil {
		global.GVA_LOG.Error("function io.Copy() Filed", zap.Any("err", copyErr.Error()))
		return "", "", errors.New("function io.Copy() Filed, err:" + copyErr.Error())
	}
	return p, filename, nil
}

func (*Local) DeleteFile(key string) error {
	p := global.GVA_CONFIG.Local.Path + "/" + key
	if strings.Contains(p, global.GVA_CONFIG.Local.Path) {
		if err := os.Remove(p); err != nil {
			return errors.New("本地文件删除失败, err:" + err.Error())
		}
	}
	return nil
}
