package p2p

import (
	"github.com/gin-gonic/gin"
	"server/service/public/p2p"
)

var svr = service.Service

func success(c *gin.Context) {
	c.JSON(200, gin.H{
		"reject":   false,
		"unchange": true,
	})
}

func reject(c *gin.Context, err string) {
	c.JSON(200, gin.H{
		"reject":        true,
		"reject_reason": err,
	})
}

// 需要修改
func change(c *gin.Context, data any) {
	c.JSON(200, gin.H{
		"unchange": false,
		"content":  data,
	})
}

// 需要修改
func unchange(c *gin.Context, data any) {
	c.JSON(200, gin.H{
		"unchange": true,
		"reject":   false,
	})
}

func New(c *gin.Context) {
	var req service.NewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(500, "解析参数失败")
		return
	}
	content, ok, err := svr.New(req)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	if !ok {
		success(c)
		return
	}
	change(c, content)
}

func Visit(c *gin.Context) {
	c.JSON(200, svr.VisitCfg(c.Query("key")))
}
