package admin

import (
	"github.com/gin-gonic/gin"
	"server/controller/admin/gost_node_rule"
	"server/global"
	"server/router/middleware"
)

func InitGostNodeRule(group *gin.RouterGroup) {
	g := group.Group("gost/node/rule", middleware.Auth(global.Jwt), middleware.AuthAdmin())
	g.POST("list", gost_node_rule.List)
}
