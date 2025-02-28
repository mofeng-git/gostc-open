package admin

import (
	"github.com/gin-gonic/gin"
	"server/controller/admin/gost_node_config"
	"server/global"
	"server/router/middleware"
)

func InitGostNodeConfig(group *gin.RouterGroup) {
	g := group.Group("gost/node/config", middleware.Auth(global.Jwt), middleware.AuthAdmin())
	g.POST("create", gost_node_config.Create)
	g.POST("delete", gost_node_config.Delete)
	g.POST("update", gost_node_config.Update)
	g.POST("page", gost_node_config.Page)
	g.POST("list", gost_node_config.List)
}
