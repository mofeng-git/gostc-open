package normal

import (
	"github.com/gin-gonic/gin"
	"server/controller/normal/gost_client_tunnel"
	"server/global"
	"server/router/middleware"
)

func InitGostClientTunnel(group *gin.RouterGroup) {
	g := group.Group("gost/client/tunnel", middleware.Auth(global.Jwt))
	g.POST("create", gost_client_tunnel.Create)
	g.POST("renew", gost_client_tunnel.Renew)
	g.POST("update", gost_client_tunnel.Update)
	g.POST("enable", gost_client_tunnel.Enable)
	g.POST("delete", gost_client_tunnel.Delete)
	g.POST("page", gost_client_tunnel.Page)
}
