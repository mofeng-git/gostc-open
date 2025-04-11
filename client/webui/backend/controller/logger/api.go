package logger

import (
	"github.com/gin-gonic/gin"
	"gostc-sub/pkg/bean"
	service "gostc-sub/webui/backend/service/logger"
)

func List(c *gin.Context) {
	bean.Response.OkData(c, service.Service.List())
}
