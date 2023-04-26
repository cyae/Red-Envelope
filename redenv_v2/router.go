package main

import (
	"github.com/gin-gonic/gin"
	"redEnv_v1/app/redEnv/handler"
)

func initRouter(r *gin.Engine) {
	r.POST("/snatch", handler.SnatchHandler)
	r.POST("/open", handler.OpenHandler)
	r.POST("/get_wallet_list", handler.GwlHandler)
}