package normal

import (
	"github.com/gin-gonic/gin"
	"server/controller/normal/gost_client_proxy"
	"server/global"
	"server/router/middleware"
)

func InitGostClientProxy(group *gin.RouterGroup) {
	g := group.Group("gost/client/proxy", middleware.Auth(global.Jwt))
	g.POST("create", gost_client_proxy.Create)
	g.POST("renew", gost_client_proxy.Renew)
	g.POST("update", gost_client_proxy.Update)
	g.POST("enable", gost_client_proxy.Enable)
	g.POST("delete", gost_client_proxy.Delete)
	g.POST("page", gost_client_proxy.Page)
}
