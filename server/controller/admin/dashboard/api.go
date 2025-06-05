package dashboard

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
	"server/service/admin/dashboard"
)

var svr = service.Service

func UserObs(c *gin.Context) {
	bean.Response.OkData(c, svr.UserObs())
}

func NodeObs(c *gin.Context) {
	bean.Response.OkData(c, svr.NodeObs())
}

func UserObsDate(c *gin.Context) {
	bean.Response.OkData(c, svr.UserObsDate(c.Query("date")))
}

func NodeObsDate(c *gin.Context) {
	bean.Response.OkData(c, svr.NodeObsDate(c.Query("date")))
}

func Count(c *gin.Context) {
	bean.Response.OkData(c, svr.Count())
}

func ClientObsDate(c *gin.Context) {
	bean.Response.OkData(c, svr.ClientObsDate(c.Query("date")))
}

func ClientHostObsDate(c *gin.Context) {
	bean.Response.OkData(c, svr.ClientHostObsDate(c.Query("date")))
}

func ClientForwardObsDate(c *gin.Context) {
	bean.Response.OkData(c, svr.ClientForwardObsDate(c.Query("date")))
}

func ClientTunnelObsDate(c *gin.Context) {
	bean.Response.OkData(c, svr.ClientTunnelObsDate(c.Query("date")))
}
