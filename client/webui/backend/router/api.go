package router

import (
	"github.com/gin-gonic/gin"
	"gostc-sub/webui/backend/controller/client"
	"gostc-sub/webui/backend/controller/config"
	"gostc-sub/webui/backend/controller/logger"
	"gostc-sub/webui/backend/controller/p2p"
	"gostc-sub/webui/backend/controller/tunnel"
)

func InitApi(engine *gin.Engine) {
	api := engine.Group("api")
	api.POST("client/create", client.Create)
	api.POST("client/update", client.Update)
	api.POST("client/delete", client.Delete)
	api.POST("client/list", client.List)
	api.POST("client/status", client.Status)

	api.POST("tunnel/create", tunnel.Create)
	api.POST("tunnel/update", tunnel.Update)
	api.POST("tunnel/delete", tunnel.Delete)
	api.POST("tunnel/list", tunnel.List)
	api.POST("tunnel/status", tunnel.Status)

	api.POST("p2p/create", p2p.Create)
	api.POST("p2p/update", p2p.Update)
	api.POST("p2p/delete", p2p.Delete)
	api.POST("p2p/list", p2p.List)
	api.POST("p2p/status", p2p.Status)

	api.POST("logger/list", logger.List)

	api.POST("config/query", config.Query)
}
