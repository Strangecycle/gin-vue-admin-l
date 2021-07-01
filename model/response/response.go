package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "操作成功", c)
}

func OkWithDetailed(data interface{}, msg string, c *gin.Context) {
	Result(SUCCESS, data, msg, c)
}

func OkWithMessage(msg string, c *gin.Context) {
	Result(SUCCESS, gin.H{}, msg, c)
}

func FailWithMessage(msg string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, msg, c)
}

func FailWithDetailed(data interface{}, msg string, c *gin.Context) {
	Result(ERROR, data, msg, c)
}
