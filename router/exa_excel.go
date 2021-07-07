package router

import (
	v1 "gin-vue-admin-l/api/v1"
	"gin-vue-admin-l/middleware"
	"github.com/gin-gonic/gin"
)

func InitExcelRouter(rg *gin.RouterGroup) {
	e := rg.Group("excel")
	e.Use(middleware.OperationRecord())
	{
		e.POST("/exportExcel", v1.ExportExcel)          // 导出 excel
		e.POST("/importExcel", v1.ImportExcel)          // 导入 excel
		e.GET("/loadExcel", v1.LoadExcel)               // 加载Excel数据
		e.GET("/downloadTemplate", v1.DownloadTemplate) // 下载模板文件
	}
}
