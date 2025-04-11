package bootstrap

import (
	"github.com/gin-gonic/gin"
	"gostc-sub/webui/backend/global"
)

var engine *gin.Engine

var Route func(e *gin.Engine)

func InitRouter() {
	gin.SetMode(gin.ReleaseMode)
	engine = gin.Default()

	if Route != nil {
		Route(engine)
	}
	global.Logger.Info("init router finish")
}
