package event

import (
	service "github.com/SianHH/frp-package/package"
	"github.com/lesismal/arpc"
)

func RemoveHandle(client *arpc.Client, callback func(key string)) {
	client.Handler.Handle("remove_config", func(c *arpc.Context) {
		var req string
		if err := c.Bind(&req); err != nil {
			_ = c.Write(err.Error())
			return
		}
		service.Del(req)
		callback(req)
		_ = c.Write("success")
	})
}
