package event

import (
	service "github.com/SianHH/frp-package/package"
	"github.com/SianHH/frp-package/package/frpc"
	v1 "github.com/SianHH/frp-package/pkg/config/v1"
	"github.com/lesismal/arpc"
	"os"
	"strconv"
	"time"
)

func HostHandle(client *arpc.Client, callback func(key, updateTag string), checkUpdate func(key, updateTag string) bool) {
	client.Handler.Handle("host_config", func(c *arpc.Context) {
		var req HostConfig
		if err := c.Bind(&req); err != nil {
			_ = c.Write(err.Error())
			return
		}

		if !checkUpdate(req.Key, req.UpdateTag) {
			_ = c.Write("success")
			return
		}

		if req.IsHttps {
			req.Http.Plugin = v1.TypedClientPluginOptions{
				Type: v1.PluginHTTP2HTTPS,
				ClientPluginOptions: &v1.HTTP2HTTPSPluginOptions{
					Type:      v1.PluginHTTP2HTTPS,
					LocalAddr: req.Http.LocalIP + ":" + strconv.Itoa(req.Http.LocalPort),
				},
			}
		}
		var proxyCfgs []v1.ProxyConfigurer
		if req.Http.Name != "" {
			proxyCfgs = append(proxyCfgs, req.Http.To())
		}
		service.Del(req.Key)
		time.Sleep(time.Second)
		req.BaseCfg.Transport.ProxyURL = os.Getenv("GOSTC_TRANSPORT_PROXY_URL")
		svc, err := frpc.NewService(req.BaseCfg, proxyCfgs, nil)
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
