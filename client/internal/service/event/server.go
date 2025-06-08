package event

import (
	service "github.com/SianHH/frp-package/package"
	"github.com/SianHH/frp-package/package/frps"
	"github.com/lesismal/arpc"
)

func ServerHandle(client *arpc.Client, httpUrl string, callback func(key string)) {
	client.Handler.Handle("server_config", func(c *arpc.Context) {
		var req ServerConfig
		if err := c.Bind(&req); err != nil {
			_ = c.Write(err.Error())
			return
		}
		service.Del(req.Key)
		for i := 0; i < len(req.ServerConfig.HTTPPlugins); i++ {
			req.ServerConfig.HTTPPlugins[i].Addr = httpUrl
		}
		svc := frps.NewService(req.ServerConfig)
		if err := svc.Start(); err != nil {
			_ = c.Write(err.Error())
			return
		}
		if err := service.Set(req.Key, svc); err != nil {
			svc.Stop()
			_ = c.Write(err.Error())
			return
		}
		callback(req.Key)
		_ = c.Write("success")
	})
}
