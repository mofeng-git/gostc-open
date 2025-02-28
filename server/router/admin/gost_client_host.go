package admin

import (
	"github.com/gin-gonic/gin"
	"server/controller/admin/gost_client_host"
	"server/global"
	"server/router/middleware"
)

func InitGostClientHost(group *gin.RouterGroup) {
	g := group.Group("gost/client/host", middleware.Auth(global.Jwt), middleware.AuthAdmin())
	g.POST("create", gost_client_host.Create)
	g.POST("config", gost_client_host.Config)
	g.POST("delete", gost_client_host.Delete)
	g.POST("page", gost_client_host.Page)
}
