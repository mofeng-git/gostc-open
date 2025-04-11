package config

import (
	"github.com/gin-gonic/gin"
	"gostc-sub/pkg/bean"
	service "gostc-sub/webui/backend/service/config"
)

func Query(c *gin.Context) {
	bean.Response.OkData(c, service.Service.Query())
}
