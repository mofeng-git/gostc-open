package system_user

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
	"server/service/normal/system_notice"
)

var svr = service.Service

func List(c *gin.Context) {
	bean.Response.OkData(c, svr.List())
}
