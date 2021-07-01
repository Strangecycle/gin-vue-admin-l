package v1

import (
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model/response"
	"gin-vue-admin-l/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Tags SysApi
// @Summary 获取所有的Api 不分页
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /api/getAllApis [post]
func GetAllApis(c *gin.Context) {
	err, apis := service.GetAllApis()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.SysAPIListResponse{Apis: apis}, "获取成功", c)
}
