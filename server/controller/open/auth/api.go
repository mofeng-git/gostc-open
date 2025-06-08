package auth

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
	"server/service/open/auth"
)

var svr = service.Service

func Checkin(c *gin.Context) {
	var req service.CheckInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	amount, err := svr.Checkin(req)
	if err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.OkData(c, amount)
}
