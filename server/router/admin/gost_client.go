package admin

import (
	"github.com/gin-gonic/gin"
	"server/controller/admin/gost_client"
	"server/global"
	"server/router/middleware"
)

func InitGostClient(group *gin.RouterGroup) {
	g := group.Group("gost/client", middleware.Auth(global.Jwt), middleware.AuthAdmin())
	g.POST("create", gost_client.Create)
	g.POST("delete", gost_client.Delete)
	g.POST("page", gost_client.Page)
	g.POST("list", gost_client.List)
}
