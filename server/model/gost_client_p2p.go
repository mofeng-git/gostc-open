package model

import (
	"fmt"
	v1 "server/pkg/p2p_cfg/v1"
	"server/pkg/utils"
)

type GostClientP2P struct {
	Base
	Name       string     `gorm:"column:name;index;comment:名称"`
	TargetIp   string     `gorm:"column:target_ip;index;comment:内网IP"`
	TargetPort string     `gorm:"column:target_port;index;comment:内网端口"`
	VKey       string     `gorm:"column:v_key;comment:访问密钥"`
	NodeCode   string     `gorm:"column:node_code;index;comment:节点编号"`
	Node       GostNode   `gorm:"foreignKey:NodeCode;references:Code"`
	ClientCode string     `gorm:"column:client_code;index;comment:客户端编号"`
	Client     GostClient `gorm:"foreignKey:ClientCode;references:Code"`
	UserCode   string     `gorm:"column:user_code;index;comment:用户编号"`
	User       SystemUser `gorm:"foreignKey:UserCode;references:Code"`
	Enable     int        `gorm:"column:enable;size:1;default:1;comment:启用状态"`
	Status     int        `gorm:"column:status;size:1;default:1;comment:状态"`
	GostClientConfig
}

func (p2p GostClientP2P) GenerateCommonCfg() (v1.ClientCommonConfig, bool) {
	if p2p.Node.P2P != 1 {
		return v1.ClientCommonConfig{}, false
	}
	return v1.ClientCommonConfig{
		ServerAddr: p2p.Node.Address,
		ServerPort: utils.StrMustInt(p2p.Node.P2PPort),
	}, true
}

func (p2p GostClientP2P) GenerateProxyCfgs() (v1.STCPProxyConfig, v1.XTCPProxyConfig) {
	if p2p.Node.P2P != 1 {
		return v1.STCPProxyConfig{}, v1.XTCPProxyConfig{}
	}
	var stcp = v1.STCPProxyConfig{
		ProxyBaseConfig: v1.ProxyBaseConfig{
			Name: "stcp_" + p2p.Code,
			Type: "stcp",
			Transport: v1.ProxyTransport{
				UseEncryption:      false,
				UseCompression:     true,
				BandwidthLimit:     fmt.Sprintf("%dKB", p2p.Limiter*128),
				BandwidthLimitMode: "client",
			},
			ProxyBackend: v1.ProxyBackend{
				LocalIP:   p2p.TargetIp,
				LocalPort: utils.StrMustInt(p2p.TargetPort),
			},
		},
		Secretkey:  p2p.VKey,
		AllowUsers: []string{"*"},
	}
	var xtcp = v1.XTCPProxyConfig{
		ProxyBaseConfig: v1.ProxyBaseConfig{
			Name: "xtcp_" + p2p.Code,
			Type: "xtcp",
			Transport: v1.ProxyTransport{
				UseEncryption:  false,
				UseCompression: true,
			},
			ProxyBackend: v1.ProxyBackend{
				LocalIP:   p2p.TargetIp,
				LocalPort: utils.StrMustInt(p2p.TargetPort),
			},
		},
		Secretkey:  p2p.VKey,
		AllowUsers: []string{"*"},
	}
	if p2p.Node.P2PDisableForward == 0 {
		return stcp, xtcp
	}
	return v1.STCPProxyConfig{}, xtcp
}

func (p2p GostClientP2P) GenerateVisitorCfgs() (v1.STCPVisitorConfig, v1.XTCPVisitorConfig) {
	if p2p.Node.P2P != 1 {
		return v1.STCPVisitorConfig{}, v1.XTCPVisitorConfig{}
	}
	return v1.STCPVisitorConfig{
			VisitorBaseConfig: v1.VisitorBaseConfig{
				Name: "v_stcp_" + p2p.Code,
				Type: "stcp",
				Transport: v1.VisitorTransport{
					UseEncryption:  false,
					UseCompression: true,
				},
				SecretKey:  p2p.VKey,
				ServerName: "stcp_" + p2p.Code,
				BindAddr:   "0.0.0.0",
			},
		},
		v1.XTCPVisitorConfig{
			VisitorBaseConfig: v1.VisitorBaseConfig{
				Name: "v_xtcp_" + p2p.Code,
				Type: "xtcp",
				Transport: v1.VisitorTransport{
					UseEncryption:  false,
					UseCompression: true,
				},
				SecretKey:  p2p.VKey,
				ServerName: "xtcp_" + p2p.Code,
			},
			KeepTunnelOpen:    true,
			MaxRetriesAnHour:  20,
			MinRetryInterval:  90,
			FallbackTo:        "v_stcp_" + p2p.Code,
			FallbackTimeoutMs: 500,
		}
}
