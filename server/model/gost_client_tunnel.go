package model

import (
	"encoding/json"
	"fmt"
	"github.com/go-gost/x/config"
	"strings"
)

type GostClientTunnel struct {
	Base
	Name       string     `gorm:"column:name;index;comment:名称"`
	TargetIp   string     `gorm:"column:target_ip;index;comment:内网IP"`
	TargetPort string     `gorm:"column:target_port;index;comment:内网端口"`
	VKey       string     `gorm:"column:v_key;comment:访问密钥"`
	NoDelay    int        `gorm:"column:no_delay;size:1;comment:无等待延迟"`
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

func (tunnel GostClientTunnel) GenerateTcpSvcConfig(chain string) (config.ServiceConfig, bool) {
	if tunnel.Node.Tunnel != 1 {
		return config.ServiceConfig{}, false
	}
	var target = tunnel.TargetIp + ":" + tunnel.TargetPort
	return config.ServiceConfig{
		Name:     "tcp_" + tunnel.Code,
		Addr:     ":0",
		Handler:  &config.HandlerConfig{Type: "rtcp", Metadata: map[string]any{"keepAlive": true}},
		Listener: &config.ListenerConfig{Type: "rtcp", Chain: chain, Metadata: map[string]any{"keepAlive": true}},
		Forwarder: &config.ForwarderConfig{
			Nodes: []*config.ForwardNodeConfig{
				{Name: target, Addr: target},
			},
		},
	}, true
}

func (tunnel GostClientTunnel) GenerateUdpSvcConfig(chain string) (config.ServiceConfig, bool) {
	if tunnel.Node.Tunnel != 1 {
		return config.ServiceConfig{}, false
	}
	var target = tunnel.TargetIp + ":" + tunnel.TargetPort
	return config.ServiceConfig{
		Name:     "udp_" + tunnel.Code,
		Addr:     ":0",
		Handler:  &config.HandlerConfig{Type: "rudp", Metadata: map[string]any{"keepAlive": true}},
		Listener: &config.ListenerConfig{Type: "rudp", Chain: chain, Metadata: map[string]any{"keepAlive": true}},
		Forwarder: &config.ForwarderConfig{
			Nodes: []*config.ForwardNodeConfig{
				{Name: target, Addr: target},
			},
		},
	}, true
}

func (tunnel GostClientTunnel) GenerateChainConfig(auth GostAuth) config.ChainConfig {
	var metadata = make(map[string]any)
	_ = json.Unmarshal([]byte(tunnel.Node.TunnelMetadata), &metadata)
	metadata["tunnel.id"] = tunnel.Code
	metadata["keepAlive"] = true

	var protocol, address string
	protocol = tunnel.Node.Protocol
	address = tunnel.Node.Address + ":" + tunnel.Node.TunnelConnPort
	replaceAddress := strings.Split(tunnel.Node.TunnelReplaceAddress, "://")
	if len(replaceAddress) == 2 {
		protocol = replaceAddress[0]
		address = replaceAddress[1]
	}
	return config.ChainConfig{
		Name: "chain_" + tunnel.Code,
		Hops: []*config.HopConfig{
			{
				Nodes: []*config.NodeConfig{
					{
						Addr: address,
						Connector: &config.ConnectorConfig{
							Type:     "tunnel",
							Metadata: metadata,
							Auth: &config.AuthConfig{
								Username: auth.User,
								Password: auth.Password,
							},
						},
						Dialer: &config.DialerConfig{
							Type: protocol,
						},
					},
				},
			},
		},
	}
}

func (tunnel GostClientTunnel) GenerateVisitTcpSvcConfig(chain string, limiter, cLimiter, rLimiter string) (config.ServiceConfig, bool) {
	return config.ServiceConfig{
		Name:     "tcp_" + tunnel.Code,
		Addr:     "",
		Limiter:  limiter,
		CLimiter: cLimiter,
		RLimiter: rLimiter,
		Handler: &config.HandlerConfig{
			Type:     "tcp",
			Metadata: map[string]any{"keepAlive": true, "nodelay": true},
			Chain:    chain,
		},
		Listener: &config.ListenerConfig{
			Type:     "tcp",
			Metadata: map[string]any{"keepAlive": true},
		},
		Forwarder: &config.ForwarderConfig{
			Nodes: []*config.ForwardNodeConfig{
				{Name: tunnel.Code, Addr: tunnel.Code},
			},
		},
		Metadata: map[string]any{"keepAlive": true},
	}, true
}

func (tunnel GostClientTunnel) GenerateVisitUdpSvcConfig(chain string, limiter, cLimiter, rLimiter string) (config.ServiceConfig, bool) {
	return config.ServiceConfig{
		Name:     "udp_" + tunnel.Code,
		Addr:     "",
		Limiter:  limiter,
		CLimiter: cLimiter,
		RLimiter: rLimiter,
		Handler: &config.HandlerConfig{
			Type:     "udp",
			Metadata: map[string]any{"keepAlive": true, "nodelay": true},
			Chain:    chain,
		},
		Listener: &config.ListenerConfig{
			Type:     "udp",
			Metadata: map[string]any{"keepAlive": true},
		},
		Forwarder: &config.ForwarderConfig{
			Nodes: []*config.ForwardNodeConfig{
				{Name: tunnel.Code, Addr: tunnel.Code},
			},
		},
		Metadata: map[string]any{"keepAlive": true},
	}, true
}

func (tunnel GostClientTunnel) GenerateVisitChainConfig(auth GostAuth) config.ChainConfig {
	var protocol, address string
	protocol = tunnel.Node.Protocol
	address = tunnel.Node.Address + ":" + tunnel.Node.TunnelConnPort
	replaceAddress := strings.Split(tunnel.Node.TunnelReplaceAddress, "://")
	if len(replaceAddress) == 2 {
		protocol = replaceAddress[0]
		address = replaceAddress[1]
	}
	return config.ChainConfig{
		Name: "chain_" + tunnel.Code,
		Hops: []*config.HopConfig{
			{
				Nodes: []*config.NodeConfig{
					{
						Addr: address,
						Connector: &config.ConnectorConfig{
							Type: "tunnel",
							Metadata: map[string]any{
								"tunnel.id": tunnel.Code,
								"nodelay":   true,
							},
							Auth: &config.AuthConfig{
								Username: auth.User,
								Password: auth.Password,
							},
						},
						Dialer: &config.DialerConfig{
							Type: protocol,
						},
					},
				},
			},
		},
	}
}

func (tunnel GostClientTunnel) GenerateVisitLimiter() config.LimiterConfig {
	return config.LimiterConfig{
		Name: "limiter_" + tunnel.Code,
		Limits: []string{
			fmt.Sprintf("$ %dKB  %dKB", tunnel.Limiter*128, tunnel.Limiter*128),
		},
	}
}

func (tunnel GostClientTunnel) GenerateVisitCLimiter() config.LimiterConfig {
	return config.LimiterConfig{
		Name: "cLimiter_" + tunnel.Code,
		Limits: []string{
			fmt.Sprintf("$ %d", tunnel.CLimiter),
		},
	}
}

func (tunnel GostClientTunnel) GenerateVisitRLimiter() config.LimiterConfig {
	return config.LimiterConfig{
		Name: "rLimiter_" + tunnel.Code,
		Limits: []string{
			fmt.Sprintf("$ %d", tunnel.RLimiter),
		},
	}
}
