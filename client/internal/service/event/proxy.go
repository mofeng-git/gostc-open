package event

import (
	service "github.com/SianHH/frp-package/package"
	"github.com/SianHH/frp-package/package/frpc"
	"github.com/SianHH/frp-package/pkg/config/types"
	v1 "github.com/SianHH/frp-package/pkg/config/v1"
	"github.com/lesismal/arpc"
)

func ProxyHandle(client *arpc.Client, callback func(key string)) {
	client.Handler.Handle("proxy_config", func(c *arpc.Context) {
		var req ProxyConfig
		if err := c.Bind(&req); err != nil {
			_ = c.Write(err.Error())
			return
		}
		var proxyCfgs []v1.ProxyConfigurer
		if req.Name != "" {
			proxyCfgs = append(proxyCfgs, &v1.TCPProxyConfig{
				ProxyBaseConfig: v1.ProxyBaseConfig{
					Name: req.Name,
					Type: "tcp",
					Transport: v1.ProxyTransport{
						UseEncryption:  true,
						UseCompression: true,
						BandwidthLimit: func() types.BandwidthQuantity {
							quantity, _ := types.NewBandwidthQuantity(req.Limiter)
							return quantity
						}(),
						BandwidthLimitMode:   "",
						ProxyProtocolVersion: "",
					},
					Metadatas: req.Metadata,
					ProxyBackend: v1.ProxyBackend{
						Plugin: v1.TypedClientPluginOptions{
							Type: "socks5",
							ClientPluginOptions: &v1.Socks5PluginOptions{
								Type:     "socks5",
								Username: req.AuthUser,
								Password: req.AuthPwd,
							},
						},
					},
				},
				RemotePort: req.Port,
			})
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
