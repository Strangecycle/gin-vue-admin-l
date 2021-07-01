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

// @Tags AuthorityMenu
// @Summary 获取用户动态路由
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body request.Empty true "空"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /menu/getBaseMenuTree [post]
func GetBaseMenuTree(c *gin.Context) {
	err, menus := service.GetBaseMenuTree()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.SysBaseMenusResponse{Menus: menus}, "获取成功", c)
}

// @Tags AuthorityMenu
// @Summary 获取指定角色menu
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetAuthorityId true "角色ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /menu/GetMenuAuthority [post]
func GetMenuAuthority(c *gin.Context) {
	var param request.GetAuthorityId
	_ = c.ShouldBindJSON(&param)
	if err := utils.Verify(param, utils.AuthorityIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err, menus := service.GetMenuAuthority(&param)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithDetailed(response.SysMenusResponse{Menus: menus}, "获取失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"menus": menus}, "获取成功", c)
}

// @Tags AuthorityMenu
// @Summary 增加menu和角色关联关系
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.AddMenuAuthorityInfo true "角色ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /menu/addMenuAuthority [post]
func AddMenuAuthority(c *gin.Context) {
	var authMenu request.AddMenuAuthorityInfo
	_ = c.ShouldBindJSON(&authMenu)
	if err := utils.Verify(authMenu, utils.AuthorityIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := service.AddMenuAuthority(authMenu.Menus, authMenu.AuthorityId)
	if err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))
		response.FailWithMessage("添加失败", c)
		return
	}
	response.OkWithMessage("添加成功", c)
}
