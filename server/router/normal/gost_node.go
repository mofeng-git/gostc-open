package normal

import (
	"github.com/gin-gonic/gin"
	"server/controller/normal/gost_node"
	"server/global"
	"server/router/middleware"
)

func InitGostNode(group *gin.RouterGroup) {
	g := group.Group("gost/node", middleware.Auth(global.Jwt))
	g.POST("list", gost_node.List)
	g.POST("cleanPort", gost_node.CleanPort)
	g.POST("create", gost_node.Create)
	g.POST("delete", gost_node.Delete)
	g.POST("update", gost_node.Update)
	g.POST("page", gost_node.Page)
}
