package event

import (
	service "github.com/SianHH/frp-package/package"
	"github.com/SianHH/frp-package/package/frpc"
	"github.com/SianHH/frp-package/package/frps"
	v1 "github.com/SianHH/frp-package/pkg/config/v1"
	"github.com/lesismal/arpc"
	"time"
)

func CustomCfgHandle(client *arpc.Client, callback func(key string)) {
	client.Handler.Handle("custom_cfg_config", func(c *arpc.Context) {
		var req CustomCfgConfig
		if err := c.Bind(&req); err != nil {
			_ = c.Write(err.Error())
			return
		}
		service.Del(req.Key)
		time.Sleep(time.Second)

		var svc service.Service
		var err error
		switch req.Type {
		case "frps":
			svc, err = frps.NewService(v1.ServerConfig{}, frps.FromBytes([]byte(req.Content)))
			if err != nil {
				_ = c.Write(err.Error())
				return
			}
		case "frpc":
			svc, err = frpc.NewService(v1.ClientCommonConfig{}, nil, nil, frpc.FromBytes([]byte(req.Content)))
			if err != nil {
				_ = c.Write(err.Error())
				return
			}
		default:
			_ = c.Write("unknown type")
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
