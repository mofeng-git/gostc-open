package event

import (
	"github.com/lesismal/arpc"
	"gostc-sub/pkg/utils"
	"strconv"
)

func PortCheckHandle(client *arpc.Client) {
	client.Handler.Handle("port_check", func(c *arpc.Context) {
		var req string
		if err := c.Bind(&req); err != nil {
			_ = c.Write(err.Error())
			return
		}
		port, err := strconv.Atoi(req)
		if err != nil {
			_ = c.Write(err.Error())
			return
		}
		if err := utils.IsUse("0.0.0.0", port); err != nil {
			_ = c.Write(err.Error())
			return
		}
		_ = c.Write("success")
	})
}
