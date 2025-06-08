package gost

import (
	"github.com/gin-gonic/gin"
	"server/service/public/gost"
)

var svr = service.Service

func Auther(c *gin.Context) {
	var req service.AutherReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}
	c.JSON(200, svr.Auther(req))
}
