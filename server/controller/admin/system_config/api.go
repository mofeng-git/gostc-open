package system_config

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
	"server/service/admin/system_config"
)

var svr = service.Service

func Base(c *gin.Context) {
	var req service.BaseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	if err := svr.Base(req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}

func Gost(c *gin.Context) {
	var req service.GostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	if err := svr.Gost(req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}

func Query(c *gin.Context) {
	var req service.QueryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	bean.Response.OkData(c, svr.Query(req))
}
