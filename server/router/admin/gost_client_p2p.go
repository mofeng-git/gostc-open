package admin

import (
	"github.com/gin-gonic/gin"
	"server/controller/admin/gost_client_p2p"
	"server/global"
	"server/router/middleware"
)

func InitGostClientP2P(group *gin.RouterGroup) {
	g := group.Group("gost/client/p2p", middleware.Auth(global.Jwt), middleware.AuthAdmin())
	g.POST("config", gost_client_p2p.Config)
	g.POST("delete", gost_client_p2p.Delete)
	g.POST("page", gost_client_p2p.Page)
}
