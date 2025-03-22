package gost_client_p2p

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
	"server/router/middleware"
	"server/service/normal/gost_client_p2p"
)

var svr = service.Service

func Create(c *gin.Context) {
	var req service.CreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	if err := svr.Create(middleware.GetClaims(c), req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}

func Renew(c *gin.Context) {
	var req service.RenewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	if err := svr.Renew(middleware.GetClaims(c), req); err != nil {
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
	if err := svr.Update(middleware.GetClaims(c), req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}

func Enable(c *gin.Context) {
	var req service.EnableReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	if err := svr.Enable(middleware.GetClaims(c), req); err != nil {
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
	if err := svr.Delete(middleware.GetClaims(c), req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}

func Page(c *gin.Context) {
	var req service.PageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	list, total := svr.Page(middleware.GetClaims(c), req)
	bean.Response.OkData(c, bean.NewPage(list, total))
}
