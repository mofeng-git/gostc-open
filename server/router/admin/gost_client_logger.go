package admin

import (
	"github.com/gin-gonic/gin"
	"server/controller/admin/gost_client_logger"
	"server/global"
	"server/router/middleware"
)

func InitGostClientLogger(group *gin.RouterGroup) {
	g := group.Group("gost/client/logger", middleware.Auth(global.Jwt))
	g.POST("page", gost_client_logger.Page)
}
