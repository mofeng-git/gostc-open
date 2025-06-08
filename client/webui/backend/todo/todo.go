package todo

import (
	"encoding/json"
	"gostc-sub/internal/common"
	service2 "gostc-sub/internal/service"
	service "gostc-sub/internal/service/visitor"
	"gostc-sub/webui/backend/bootstrap"
	"gostc-sub/webui/backend/global"
	"gostc-sub/webui/backend/model"
	"strconv"
)

func init() {
	bootstrap.TodoFunc = func() {
		for _, key := range global.ClientFS.ListKeys() {
			value, ok := global.ClientFS.Get(key)
			if !ok {
				continue
			}
			var client model.Client
			marshal, _ := json.Marshal(value)
			_ = json.Unmarshal(marshal, &client)
			svc := service2.NewClient(common.GenerateWsUrl(client.Tls == 1, client.Address), common.GenerateHttpUrl(client.Tls == 1, client.Address), client.Key)
			global.ClientMap.Store(client.Key, svc)
			if client.AutoStart == 1 {
				_ = svc.Start()
			}
		}

		for _, key := range global.P2PFS.ListKeys() {
			value, ok := global.P2PFS.Get(key)
			if !ok {
				continue
			}
			var p2p model.P2P
			marshal, _ := json.Marshal(value)
			_ = json.Unmarshal(marshal, &p2p)
			port, _ := strconv.Atoi(p2p.Port)
			svc := service.NewP2P(common.GenerateHttpUrl(p2p.Tls == 1, p2p.Address), p2p.Key, p2p.Bind, port)
			global.P2PMap.Store(p2p.Key, svc)
			if p2p.AutoStart == 1 {
				_ = svc.Start()
			}
		}

		for _, key := range global.TunnelFS.ListKeys() {
			value, ok := global.TunnelFS.Get(key)
			if !ok {
				continue
			}
			var tunnel model.Tunnel
			marshal, _ := json.Marshal(value)
			_ = json.Unmarshal(marshal, &tunnel)
			port, _ := strconv.Atoi(tunnel.Port)
			svc := service.NewTunnel(common.GenerateHttpUrl(tunnel.Tls == 1, tunnel.Address), tunnel.Key, tunnel.Bind, port)
			global.TunnelMap.Store(tunnel.Key, svc)
			if tunnel.AutoStart == 1 {
				_ = svc.Start()
			}
		}
	}
}
