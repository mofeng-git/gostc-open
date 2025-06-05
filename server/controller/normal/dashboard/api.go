package dashboard

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
	"server/router/middleware"
	"server/service/normal/dashboard"
)

var svr = service.Service

func Count(c *gin.Context) {
	claims := middleware.GetClaims(c)
	bean.Response.OkData(c, svr.Count(claims))
}

func ClientObsDate(c *gin.Context) {
	claims := middleware.GetClaims(c)
	bean.Response.OkData(c, svr.ClientObsDate(claims, c.Query("date")))
}

func ClientHostObsDate(c *gin.Context) {
	claims := middleware.GetClaims(c)
	bean.Response.OkData(c, svr.ClientHostObsDate(claims, c.Query("date")))
}

func ClientForwardObsDate(c *gin.Context) {
	claims := middleware.GetClaims(c)
	bean.Response.OkData(c, svr.ClientForwardObsDate(claims, c.Query("date")))
}

func ClientTunnelObsDate(c *gin.Context) {
	claims := middleware.GetClaims(c)
	bean.Response.OkData(c, svr.ClientTunnelObsDate(claims, c.Query("date")))
}
