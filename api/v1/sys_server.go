package v1

import (
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model/response"
	"gin-vue-admin-l/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Tags System
// @Summary 获取服务器信息
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /system/getServerInfo [post]
func GetServerInfo(c *gin.Context) {
	server, err := service.GetServerInfo()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"server": server}, "获取成功", c)
}
