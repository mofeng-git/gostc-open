package gost

import (
	"github.com/gin-gonic/gin"
	"server/service/public/gost"
)

var svr = service.Service

func ClientWs(c *gin.Context) {
	svr.ClientWs(c)
}

func NodeWs(c *gin.Context) {
	svr.NodeWs(c)
}

func NodePort(c *gin.Context) {
	var req service.NodePortReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}
	svr.NodePort(req)
}

func ClientPort(c *gin.Context) {
	var req service.ClientPortReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}
	svr.ClientPort(req)
}

func Ingress(c *gin.Context) {
	c.JSON(200, svr.Ingress(c.Query("host")))
}

func Auther(c *gin.Context) {
	var req service.AutherReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}
	c.JSON(200, svr.Auther(req))
}

func Admission(c *gin.Context) {
	var req service.AdmissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}
	c.JSON(200, svr.Admission(req))
}

func Limiter(c *gin.Context) {
	var req service.LimiterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}
	c.JSON(200, svr.Limiter(req))
}

func Obs(c *gin.Context) {
	var req service.ObsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}
	c.JSON(200, svr.Obs(c.Query("tunnel"), req))
}

func Visit(c *gin.Context) {
	c.JSON(200, svr.VisitCfg(c.Query("key")))
}
