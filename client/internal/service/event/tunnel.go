package event

import (
	service "github.com/SianHH/frp-package/package"
	"github.com/SianHH/frp-package/package/frpc"
	v1 "github.com/SianHH/frp-package/pkg/config/v1"
	"github.com/lesismal/arpc"
)

func TunnelHandle(client *arpc.Client, callback func(key string)) {
	client.Handler.Handle("tunnel_config", func(c *arpc.Context) {
		var req TunnelConfig
		if err := c.Bind(&req); err != nil {
			_ = c.Write(err.Error())
			return
		}
		var proxyCfgs []v1.ProxyConfigurer
		if req.STCP.Name != "" {
			proxyCfgs = append(proxyCfgs, req.STCP.To())
		}
		if req.SUDP.Name != "" {
			proxyCfgs = append(proxyCfgs, req.SUDP.To())
		}
		service.Del(req.Key)
		svc := frpc.NewService(req.BaseCfg, proxyCfgs, nil)
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
