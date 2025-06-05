package admin

import (
	"github.com/gin-gonic/gin"
	"server/controller/admin/gost_client_tunnel"
	"server/global"
	"server/router/middleware"
)

func InitGostClientTunnel(group *gin.RouterGroup) {
	g := group.Group("gost/client/tunnel", middleware.Auth(global.Jwt), middleware.AuthAdmin())
	g.POST("config", gost_client_tunnel.Config)
	g.POST("delete", gost_client_tunnel.Delete)
	g.POST("page", gost_client_tunnel.Page)
}
