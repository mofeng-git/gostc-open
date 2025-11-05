package admin

import (
	"github.com/gin-gonic/gin"
	"server/controller/admin/system_config"
	"server/global"
	"server/router/middleware"
)

func InitSystemConfig(group *gin.RouterGroup) {
    g := group.Group("system/config", middleware.Auth(global.Jwt), middleware.AuthAdmin())
    g.POST("base", system_config.Base)
    g.POST("gost", system_config.Gost)
    g.POST("email", system_config.Email)
    g.POST("emailVerify", system_config.EmailVerify)
    g.POST("query", system_config.Query)
    g.POST("home", system_config.Home)
}
