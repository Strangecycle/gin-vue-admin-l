package request

import "gin-vue-admin-l/model"

type SysOperationRecordSearch struct {
	model.SysOperationRecord
	PageInfo
}
