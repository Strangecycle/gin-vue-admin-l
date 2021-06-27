package v1

import (
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model/response"
	"gin-vue-admin-l/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Tags AuthorityMenu
// @Summary 获取用户动态路由
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body request.Empty true "空"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /menu/getMenu [post]
func GetMenu(c *gin.Context) {
	err, menus := service.GetMenuTree(getUserAuthorityId(c))
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.SysMenusResponse{Menus: menus}, "获取成功", c)
}
