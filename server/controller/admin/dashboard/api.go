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

func Count(c *gin.Context) {
	bean.Response.OkData(c, svr.Count())
}
