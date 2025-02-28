package admin

import (
	"github.com/gin-gonic/gin"
	"server/controller/admin/system_user"
	"server/global"
	"server/router/middleware"
)

func InitSystemUser(group *gin.RouterGroup) {
	g := group.Group("system/user", middleware.Auth(global.Jwt), middleware.AuthAdmin())
	g.POST("create", system_user.Create)
	g.POST("update", system_user.Update)
	g.POST("delete", system_user.Delete)
	g.POST("page", system_user.Page)
	g.POST("list", system_user.List)
}
