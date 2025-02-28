package gost

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
	"server/router/middleware"
	"server/service/normal/gost"
)

var svr = service.Service

func Info(c *gin.Context) {
	bean.Response.OkData(c, svr.Info(middleware.GetClaims(c)))
}
