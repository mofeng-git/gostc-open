package model

import (
	"fmt"
	"github.com/go-gost/x/config"
	"strings"
)

type GostClientProxy struct {
	Base
	Name       string     `gorm:"column:name;index;comment:名称"`
	Protocol   string     `gorm:"column:protocol;index;comment:代理类型"`
	Port       string     `gorm:"column:port;comment:访问端口"`
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

func (proxy *GostClientProxy) GenerateSvcConfig(chain, limiter, cLimiter, rLimiter, obs string) (clientCfg config.ServiceConfig, ok bool) {
	if proxy.Node.Proxy != 1 {
		return clientCfg, ok
	}

	clientCfg = config.ServiceConfig{
		Name:     "proxy_" + proxy.Code,
		Addr:     ":" + proxy.Port,
		Limiter:  limiter,
		CLimiter: cLimiter,
		RLimiter: rLimiter,
		Observer: obs,
		Handler:  &config.HandlerConfig{Type: proxy.Protocol, Chain: chain},
		Listener: &config.ListenerConfig{Type: "tcp"},
		Metadata: map[string]any{
			"enableStats":           true,
			"observer.period":       "60s",
			"observer.resetTraffic": true,
		},
	}
	return clientCfg, true
}

func (proxy *GostClientProxy) GenerateChainConfig(auth GostAuth) config.ChainConfig {
	var protocol, address string
	protocol = proxy.Node.Protocol
	address = proxy.Node.Address + ":" + proxy.Node.ForwardConnPort
	replaceAddress := strings.Split(proxy.Node.ForwardReplaceAddress, "://")
	if len(replaceAddress) == 2 {
		protocol = replaceAddress[0]
		address = replaceAddress[1]
	}

	return config.ChainConfig{
		Name: "chain_" + proxy.Code,
		Hops: []*config.HopConfig{
			{
				Nodes: []*config.NodeConfig{
					{
						Addr: address,
						Connector: &config.ConnectorConfig{
							Type: "relay",
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

func (proxy *GostClientProxy) GenerateLimiter() config.LimiterConfig {
	return config.LimiterConfig{
		Name: "limiter_" + proxy.Code,
		Limits: []string{
			fmt.Sprintf("$ %dKB  %dKB", proxy.Limiter*128, proxy.Limiter*128),
		},
	}
}

func (proxy *GostClientProxy) GenerateCLimiter() config.LimiterConfig {
	return config.LimiterConfig{
		Name: "cLimiter_" + proxy.Code,
		Limits: []string{
			fmt.Sprintf("$ %d", proxy.CLimiter),
		},
	}
}

func (proxy *GostClientProxy) GenerateRLimiter() config.LimiterConfig {
	return config.LimiterConfig{
		Name: "rLimiter_" + proxy.Code,
		Limits: []string{
			fmt.Sprintf("$ %d", proxy.RLimiter),
		},
	}
}

func (proxy *GostClientProxy) GenerateObs(host, nodeVersion string) config.ObserverConfig {
	if nodeVersion > "v1.1.2" {
		return config.ObserverConfig{}
	}
	return config.ObserverConfig{
		Name: "obs_" + proxy.Code,
		Plugin: &config.PluginConfig{
			Type: "http",
			Addr: fmt.Sprintf("%s/api/v1/public/gost/obs?tunnel=%s", host, proxy.Code),
		},
	}
}
