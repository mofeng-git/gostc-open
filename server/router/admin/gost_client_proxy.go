package admin

import (
	"github.com/gin-gonic/gin"
	"server/controller/admin/gost_client_proxy"
	"server/global"
	"server/router/middleware"
)

func InitGostClientProxy(group *gin.RouterGroup) {
	g := group.Group("gost/client/proxy", middleware.Auth(global.Jwt), middleware.AuthAdmin())
	g.POST("config", gost_client_proxy.Config)
	g.POST("delete", gost_client_proxy.Delete)
	g.POST("page", gost_client_proxy.Page)
}
