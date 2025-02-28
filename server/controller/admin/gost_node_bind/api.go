package gost_node

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
	"server/service/admin/gost_node_bind"
)

var svr = service.Service

func Update(c *gin.Context) {
	var req service.UpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	if err := svr.Update(req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}
