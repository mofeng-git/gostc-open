package admin

import (
	"github.com/gin-gonic/gin"
	"server/controller/admin/gost_client_forward"
	"server/global"
	"server/router/middleware"
)

func InitGostClientForward(group *gin.RouterGroup) {
	g := group.Group("gost/client/forward", middleware.Auth(global.Jwt), middleware.AuthAdmin())
	g.POST("create", gost_client_forward.Create)
	g.POST("config", gost_client_forward.Config)
	g.POST("delete", gost_client_forward.Delete)
	g.POST("page", gost_client_forward.Page)
}
