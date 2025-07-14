package event

import (
	service "github.com/SianHH/frp-package/package"
	"github.com/SianHH/frp-package/package/frpc"
	v1 "github.com/SianHH/frp-package/pkg/config/v1"
	"github.com/lesismal/arpc"
	"os"
)

func ForwardHandle(client *arpc.Client, callback func(key string)) {
	client.Handler.Handle("forward_config", func(c *arpc.Context) {
		var req ForwardConfig
		if err := c.Bind(&req); err != nil {
			_ = c.Write(err.Error())
			return
		}
		var proxyCfgs []v1.ProxyConfigurer
		if req.TCP.Name != "" {
			proxyCfgs = append(proxyCfgs, req.TCP.To())
		}
		if req.UDP.Name != "" {
			proxyCfgs = append(proxyCfgs, req.UDP.To())
		}
		service.Del(req.Key)
		req.BaseCfg.Transport.ProxyURL = os.Getenv("GOSTC_TRANSPORT_PROXY_URL")
		svc, err := frpc.NewService(req.BaseCfg, proxyCfgs, nil)
		if err != nil {
			_ = c.Write(err.Error())
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
		callback(req.Key)
		_ = c.Write("success")
	})
}
