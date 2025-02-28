package normal

import (
	"github.com/gin-gonic/gin"
	"server/controller/normal/gost"
	"server/global"
	"server/router/middleware"
)

func InitGost(group *gin.RouterGroup) {
	g := group.Group("gost", middleware.Auth(global.Jwt))
	g.POST("info", gost.Info)
}
