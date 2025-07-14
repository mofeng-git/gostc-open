package event

import (
	service "github.com/SianHH/frp-package/package"
	"github.com/SianHH/frp-package/package/frpc"
	v1 "github.com/SianHH/frp-package/pkg/config/v1"
	"github.com/lesismal/arpc"
	"os"
)

func P2PHandle(client *arpc.Client, callback func(key string)) {
	client.Handler.Handle("p2p_config", func(c *arpc.Context) {
		var req P2PConfig
		if err := c.Bind(&req); err != nil {
			_ = c.Write(err.Error())
			return
		}
		var proxyCfgs []v1.ProxyConfigurer
		if req.XTCP.Name != "" {
			proxyCfgs = append(proxyCfgs, req.XTCP.To())
		}
		if req.STCP.Name != "" {
			proxyCfgs = append(proxyCfgs, req.STCP.To())
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
