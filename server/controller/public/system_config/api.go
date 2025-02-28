package system_config

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
	"server/service/public/system_config"
)

var svr = service.Service

func Query(c *gin.Context) {
	bean.Response.OkData(c, svr.Query())
}
