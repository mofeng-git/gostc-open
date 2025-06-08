package event

import (
	"github.com/lesismal/arpc"
	"gostc-sub/internal/common"
)

func StopHandle(client *arpc.Client, callback func()) {
	client.Handler.Handle("stop", func(c *arpc.Context) {
		var req string
		if err := c.Bind(&req); err != nil {
			_ = c.Write(err.Error())
			return
		}
		common.Logger.AddLog("Client", "停止原因："+req)
		client.Stop()
		callback()
		_ = c.Write("success")
	})
}
