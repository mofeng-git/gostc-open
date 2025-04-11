package router

import (
	"github.com/gin-gonic/gin"
	"gostc-sub/webui/backend/bootstrap"
)

func init() {
	bootstrap.Route = func(engine *gin.Engine) {
		InitStatic(engine)
		InitApi(engine)
	}
}
