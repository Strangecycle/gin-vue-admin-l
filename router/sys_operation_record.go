package router

import (
	v1 "gin-vue-admin-l/api/v1"
	"gin-vue-admin-l/middleware"
	"github.com/gin-gonic/gin"
)

func InitSysOperationRecordRouter(rg *gin.RouterGroup) {
	s := rg.Group("sysOperationRecord")
	s.Use(middleware.OperationRecord())
	{
		s.GET("getSysOperationRecordList", v1.GetSysOperationRecordList)            // 获取 SysOperationRecord 列表
		s.DELETE("deleteSysOperationRecord", v1.DeleteSysOperationRecord)           // 删除 SysOperationRecord
		s.DELETE("deleteSysOperationRecordByIds", v1.DeleteSysOperationRecordByIds) // 批量删除SysOperationRecord
	}
}
