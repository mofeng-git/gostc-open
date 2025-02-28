package admin

import (
	"github.com/gin-gonic/gin"
	"server/controller/admin/gost_node_logger"
	"server/global"
	"server/router/middleware"
)

func InitGostNodeLogger(group *gin.RouterGroup) {
	g := group.Group("gost/node/logger", middleware.Auth(global.Jwt))
	g.POST("page", gost_node_logger.Page)
}
