package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func initServer(addr string, r *gin.Engine) server {
	// 实例化一个 http 实例
	return &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
