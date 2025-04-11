package p2p

import (
	"github.com/gin-gonic/gin"
	"gostc-sub/pkg/bean"
	service "gostc-sub/webui/backend/service/p2p"
)

func Create(c *gin.Context) {
	var req service.CreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	if err := service.Service.Create(req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}

func Delete(c *gin.Context) {
	var req service.DeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	if err := service.Service.Delete(req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}

func Update(c *gin.Context) {
	var req service.UpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	if err := service.Service.Update(req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}

func List(c *gin.Context) {
	bean.Response.OkData(c, service.Service.List())
}

func Status(c *gin.Context) {
	var req service.StatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	if err := service.Service.Status(req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}
