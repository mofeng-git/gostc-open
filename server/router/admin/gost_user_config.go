package admin

import (
	"github.com/gin-gonic/gin"
	"server/controller/admin/gost_user_config"
	"server/global"
	"server/router/middleware"
)

func InitGostUserConfig(group *gin.RouterGroup) {
	g := group.Group("gost/user/config", middleware.Auth(global.Jwt), middleware.AuthAdmin())
	g.POST("create", gost_user_config.Create)
	g.POST("delete", gost_user_config.Delete)
	g.POST("update", gost_user_config.Update)
	g.POST("page", gost_user_config.Page)
	g.POST("list", gost_user_config.List)
	g.POST("node/list", gost_user_config.NodeList)
}
