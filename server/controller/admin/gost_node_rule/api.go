package gost_node_rule

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
	"server/service/admin/gost_node_rule"
)

var svr = service.Service

func List(c *gin.Context) {
	bean.Response.OkData(c, svr.List())
}
