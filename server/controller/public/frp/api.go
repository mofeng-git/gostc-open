package frp

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"server/service/public/frp"
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

func Login(c *gin.Context) {
	var req service.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(500, "解析参数失败")
		return
	}
	content, err := svr.Login(req)
	if err != nil {
		reject(c, err.Error())
		return
	}
	if content == nil {
		success(c)
	} else {
		change(c, content)
	}
}

func NewProxy(c *gin.Context) {
	var req service.NewProxyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(500, "解析参数失败")
		return
	}
	content, err := svr.NewProxy(req)
	if err != nil {
		reject(c, err.Error())
		return
	}
	if content == nil {
		success(c)
	} else {
		change(c, content)
	}
}

// {"version":"0.1.0","op":"CloseProxy","content":{"user":{"user":"","metas":{"password":"2l5y63jh3r","user":"5gzxwh2ntg"},"run_id":"202d484d1db5f6ef"},"proxy_name":"bfe50ee4-72b7-495e-84a3-fe43b0a49321_proxy"}}
type CloseProxyReq struct {
	Version string `json:"version"`
	Op      string `json:"op"`
	Content struct {
		User struct {
			User  string `json:"user"`
			Metas struct {
				Password string `json:"password"`
				User     string `json:"user"`
			} `json:"metas"`
			RunId string `json:"run_id"`
		} `json:"user"`
		ProxyName string `json:"proxy_name"`
	} `json:"content"`
}

func CloseProxy(c *gin.Context) {
	success(c)
}

func Ping(c *gin.Context) {
	success(c)
}

func NewWorkConn(c *gin.Context) {
	var req service.NewWorkConnReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(500, "解析参数失败")
		return
	}
	content, err := svr.NewWorkConn(req)
	if err != nil {
		reject(c, err.Error())
		return
	}
	if content == nil {
		success(c)
	} else {
		change(c, content)
	}
}

func NewUserConn(c *gin.Context) {
	var req service.NewUserConnReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(500, "解析参数失败")
		return
	}
	content, err := svr.NewUserConn(req)
	if err != nil {
		reject(c, err.Error())
		return
	}
	if content == nil {
		success(c)
	} else {
		change(c, content)
	}
}

func VisitorTunnel(c *gin.Context) {
	key := c.Query("key")
	result, err := svr.VisitorTunnel(key)
	if err != nil {
		_ = c.Error(err)
		return
	}
	marshal, _ := yaml.Marshal(result)
	_, _ = c.Writer.Write(marshal)
}

func VisitorP2P(c *gin.Context) {
	key := c.Query("key")
	result, err := svr.VisitorP2P(key)
	if err != nil {
		_ = c.Error(err)
		return
	}
	marshal, _ := yaml.Marshal(result)
	_, _ = c.Writer.Write(marshal)
}
