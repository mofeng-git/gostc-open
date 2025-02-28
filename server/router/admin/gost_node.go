package admin

import (
	"github.com/gin-gonic/gin"
	"server/controller/admin/gost_node"
	"server/global"
	"server/router/middleware"
)

func InitGostNode(group *gin.RouterGroup) {
	g := group.Group("gost/node", middleware.Auth(global.Jwt), middleware.AuthAdmin())
	g.POST("create", gost_node.Create)
	g.POST("delete", gost_node.Delete)
	g.POST("update", gost_node.Update)
	g.POST("query", gost_node.Query)
	g.POST("page", gost_node.Page)
	g.POST("list", gost_node.List)
}
