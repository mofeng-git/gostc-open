package common

import (
	"gostc-sub/p2p/frpc"
	v1 "gostc-sub/p2p/pkg/config/v1"
	registry2 "gostc-sub/p2p/registry"
)

type ClientP2PConfigData struct {
	Code    string
	Common  v1.ClientCommonConfig
	STCPCfg v1.STCPProxyConfig
	XTCPCfg v1.XTCPProxyConfig
}

func WsP2PConfig(data ClientP2PConfigData) {
	registry2.Del(data.Code)
	svc := frpc.NewService(data.Common, []v1.ProxyConfigurer{
		&data.STCPCfg,
		&data.XTCPCfg,
	}, nil)
	if err := svc.Start(); err == nil {
		_ = registry2.Set(data.Code, svc)
		SvcMap[data.Code] = true
	}
}
