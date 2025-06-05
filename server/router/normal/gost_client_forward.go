package normal

import (
	"github.com/gin-gonic/gin"
	"server/controller/normal/gost_client_forward"
	"server/global"
	"server/router/middleware"
)

func InitGostClientForward(group *gin.RouterGroup) {
	g := group.Group("gost/client/forward", middleware.Auth(global.Jwt))
	g.POST("create", gost_client_forward.Create)
	g.POST("admission", gost_client_forward.Admission)
	g.POST("renew", gost_client_forward.Renew)
	g.POST("update", gost_client_forward.Update)
	g.POST("migrate", gost_client_forward.Migrate)
	g.POST("matcher", gost_client_forward.Matcher)
	g.POST("enable", gost_client_forward.Enable)
	g.POST("delete", gost_client_forward.Delete)
	g.POST("page", gost_client_forward.Page)
}
