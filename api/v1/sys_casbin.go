package v1

import (
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model/request"
	"gin-vue-admin-l/model/response"
	"gin-vue-admin-l/service"
	"gin-vue-admin-l/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Tags Casbin
// @Summary 获取权限列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CasbinInReceive true "权限id, 权限模型列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /casbin/getPolicyPathByAuthorityId [post]
func GetPolicyPathByAuthorityId(c *gin.Context) {
	var casbin request.CasbinInReceive
	_ = c.ShouldBindJSON(&casbin)
	if err := utils.Verify(casbin, utils.AuthorityIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	paths := service.GetPolicyPathByAuthorityId(casbin.AuthorityId)
	response.OkWithDetailed(response.PolicyPathResponse{Paths: paths}, "获取成功", c)
}

// @Tags Casbin
// @Summary 更新角色api权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CasbinInReceive true "权限id, 权限模型列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /casbin/UpdateCasbin [post]
func UpdateCasbin(c *gin.Context) {
	var cir request.CasbinInReceive
	_ = c.ShouldBindJSON(&cir)
	if err := utils.Verify(cir, utils.AuthorityIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := service.UpdateCasbin(cir.AuthorityId, cir.CasbinInfos)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
