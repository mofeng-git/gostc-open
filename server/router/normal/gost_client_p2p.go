package normal

import (
	"github.com/gin-gonic/gin"
	"server/controller/normal/gost_client_p2p"
	"server/global"
	"server/router/middleware"
)

func InitGostClientP2P(group *gin.RouterGroup) {
	g := group.Group("gost/client/p2p", middleware.Auth(global.Jwt))
	g.POST("create", gost_client_p2p.Create)
	g.POST("renew", gost_client_p2p.Renew)
	g.POST("migrate", gost_client_p2p.Migrate)
	g.POST("update", gost_client_p2p.Update)
	g.POST("enable", gost_client_p2p.Enable)
	g.POST("delete", gost_client_p2p.Delete)
	g.POST("page", gost_client_p2p.Page)
}
