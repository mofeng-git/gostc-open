package normal

import (
	"github.com/gin-gonic/gin"
	"server/controller/normal/gost_client_host"
	"server/global"
	"server/router/middleware"
)

func InitGostClientHost(group *gin.RouterGroup) {
	g := group.Group("gost/client/host", middleware.Auth(global.Jwt))
	g.POST("create", gost_client_host.Create)
	g.POST("admission", gost_client_host.Admission)
	g.POST("renew", gost_client_host.Renew)
	g.POST("update", gost_client_host.Update)
	g.POST("enable", gost_client_host.Enable)
	g.POST("delete", gost_client_host.Delete)
	g.POST("page", gost_client_host.Page)
}
