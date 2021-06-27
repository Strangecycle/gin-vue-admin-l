package router

import (
	v1 "gin-vue-admin-l/api/v1"
	"gin-vue-admin-l/middleware"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(rg *gin.RouterGroup) {
	u := rg.Group("user")
	u.Use(middleware.OperationRecord())
	{
		u.POST("getUserList", v1.GetUserList)           // 分页获取用户列表
		u.POST("register", v1.Register)                 // 注册用户
		u.DELETE("deleteUser", v1.DeleteUser)           // 删除用户
		u.POST("setUserAuthority", v1.SetUserAuthority) // 设置用户角色
		u.POST("changePassword", v1.ChangePassword)     // 修改密码
		// TODO 设置用户信息
	}
}
