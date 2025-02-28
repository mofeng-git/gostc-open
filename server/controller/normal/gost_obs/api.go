package gost_obs

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
	"server/router/middleware"
	service "server/service/normal/gost_obs"
)

var svr = service.Service

func TunnelMonth(c *gin.Context) {
	var req service.TunnelMonthReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	bean.Response.OkData(c, svr.TunnelMonth(req))
}

func ClientMonth(c *gin.Context) {
	var req service.ClientMonthReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	bean.Response.OkData(c, svr.ClientMonth(req))
}

func NodeMonth(c *gin.Context) {
	var req service.NodeMonthReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	bean.Response.OkData(c, svr.NodeMonth(req))
}

func UserMonth(c *gin.Context) {
	var req service.UserMonthReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	bean.Response.OkData(c, svr.UserMonth(middleware.GetClaims(c), req))
}
