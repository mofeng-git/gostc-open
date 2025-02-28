package normal

import (
	"github.com/gin-gonic/gin"
	"server/controller/normal/system_notice"
	"server/global"
	"server/router/middleware"
)

func InitSystemNotice(group *gin.RouterGroup) {
	g := group.Group("system/notice", middleware.Auth(global.Jwt))
	g.POST("list", system_user.List)
}
