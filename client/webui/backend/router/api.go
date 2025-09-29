package router

import (
	"github.com/gin-gonic/gin"
	"gostc-sub/webui/backend/controller/client"
	"gostc-sub/webui/backend/controller/config"
	"gostc-sub/webui/backend/controller/logger"
)

func InitApi(engine *gin.Engine) {
	// 注册带前缀的API路由
	api := engine.Group("/extras/gostc/api")
	api.POST("client/create", client.Create)
	api.POST("client/update", client.Update)
	api.POST("client/delete", client.Delete)
	api.POST("client/list", client.List)
	api.POST("client/status", client.Status)

	api.POST("logger/list", logger.List)

	api.POST("config/query", config.Query)
}
