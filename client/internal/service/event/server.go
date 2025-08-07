package event

import (
	service "github.com/SianHH/frp-package/package"
	"github.com/SianHH/frp-package/package/frps"
	"github.com/lesismal/arpc"
	"time"
)

func ServerHandle(client *arpc.Client, httpUrl string, callback func(key, updateTag string), checkUpdate func(key, updateTag string) bool) {
	client.Handler.Handle("server_config", func(c *arpc.Context) {
		var req ServerConfig
		if err := c.Bind(&req); err != nil {
			_ = c.Write(err.Error())
			return
		}

		if !checkUpdate(req.Key, req.UpdateTag) {
			_ = c.Write("success")
			return
		}

		service.Del(req.Key)
		time.Sleep(time.Second)
		for i := 0; i < len(req.ServerConfig.HTTPPlugins); i++ {
			req.ServerConfig.HTTPPlugins[i].Addr = httpUrl
		}
		svc, err := frps.NewService(req.ServerConfig)
		if err != nil {
			_ = c.Write(err)
			return
		}
		if err := svc.Start(); err != nil {
			_ = c.Write(err.Error())
			return
		}
		if err := service.Set(req.Key, svc); err != nil {
			svc.Stop()
			_ = c.Write(err.Error())
			return
		}
		callback(req.Key, req.UpdateTag)
		_ = c.Write("success")
	})
}
