package response

import "gin-vue-admin-l/model"

type ExaFileResponse struct {
	File model.ExaFileUploadAndDownload `json:"file"`
}
