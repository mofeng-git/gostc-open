package engine

import (
	"gostc-sub/pkg/p2p/frpc"
	v1 "gostc-sub/pkg/p2p/pkg/config/v1"
	"gostc-sub/pkg/p2p/registry"
)

type ClientP2PConfigData struct {
	Code    string
	Common  v1.ClientCommonConfig
	STCPCfg v1.STCPProxyConfig
	XTCPCfg v1.XTCPProxyConfig
}

func (e *Event) WsP2PConfig(data ClientP2PConfigData) {
	registry.Del(data.Code)
	svc := frpc.NewService(data.Common, []v1.ProxyConfigurer{
		&data.STCPCfg,
		&data.XTCPCfg,
	}, nil)
	if err := svc.Start(); err == nil {
		_ = registry.Set(data.Code, svc)
		e.svcMap.Store(data.Code, true)
	}
}
