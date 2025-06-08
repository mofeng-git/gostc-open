package service

import v1 "github.com/SianHH/frp-package/pkg/config/v1"

type TunnelConfig struct {
	Common v1.ClientCommonConfig
	STCP   v1.STCPVisitorConfig
	SUDP   v1.SUDPVisitorConfig
}

type P2PConfig struct {
	Common v1.ClientCommonConfig
	XTCP   v1.XTCPVisitorConfig
	STCP   v1.STCPVisitorConfig
}
