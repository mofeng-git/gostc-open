package admin

import (
	"github.com/gin-gonic/gin"
	"server/controller/admin/gost_node_bind"
	"server/global"
	"server/router/middleware"
)

func InitGostNodeBind(group *gin.RouterGroup) {
	g := group.Group("gost/node/bind", middleware.Auth(global.Jwt), middleware.AuthAdmin())
	g.POST("update", gost_node.Update)
}
