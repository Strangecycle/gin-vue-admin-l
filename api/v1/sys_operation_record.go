package v1

import (
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model"
	"gin-vue-admin-l/model/request"
	"gin-vue-admin-l/model/response"
	"gin-vue-admin-l/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Tags SysOperationRecord
// @Summary 分页获取SysOperationRecord列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysOperationRecordSearch true "页码, 每页大小, 搜索条件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sysOperationRecord/getSysOperationRecordList [get]
func GetSysOperationRecordList(c *gin.Context) {
	var pageInfo request.SysOperationRecordSearch
	_ = c.ShouldBindQuery(&pageInfo)

	err, list, total := service.GetSysOperationRecordInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// @Tags SysOperationRecord
// @Summary 删除SysOperationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysOperationRecord true "SysOperationRecord模型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysOperationRecord/deleteSysOperationRecord [delete]
func DeleteSysOperationRecord(c *gin.Context) {
	var record model.SysOperationRecord
	_ = c.ShouldBindJSON(&record)

	err := service.DeleteSysOperationRecord(record)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// @Tags SysOperationRecord
// @Summary 批量删除SysOperationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除SysOperationRecord"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /sysOperationRecord/deleteSysOperationRecordByIds [delete]
func DeleteSysOperationRecordByIds(c *gin.Context) {
	var ids request.IdsReq
	_ = c.ShouldBindJSON(&ids)
	err := service.DeleteSysOperationRecordByIds(ids)

	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage("批量删除失败", c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}
