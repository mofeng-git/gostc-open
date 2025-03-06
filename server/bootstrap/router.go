package bootstrap

import (
	"github.com/gin-gonic/gin"
	"server/global"
)

var engine *gin.Engine

var Route func(e *gin.Engine)

func InitRouter() {
	if global.MODE == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	engine = gin.Default()
	
	if Route != nil {
		Route(engine)
	}
	global.Logger.Info("init router finish")
}
