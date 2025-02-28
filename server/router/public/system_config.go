package public

import (
	"github.com/gin-gonic/gin"
	"server/controller/public/system_config"
)

func InitSystemConfig(group *gin.RouterGroup) {
	g := group.Group("system/config")
	g.POST("query", system_config.Query)
}
