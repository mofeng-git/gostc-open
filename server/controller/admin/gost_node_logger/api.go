package gost_node_logger

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
	service "server/service/admin/gost_node_logger"
)

var svr = service.Service

func Page(c *gin.Context) {
	var req service.PageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	list, total := svr.Page(req)
	bean.Response.OkData(c, bean.NewPage(list, total))
}
