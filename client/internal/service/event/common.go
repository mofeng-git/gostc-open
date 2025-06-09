package event

import (
	"github.com/SianHH/frp-package/pkg/config/types"
	v1 "github.com/SianHH/frp-package/pkg/config/v1"
	"github.com/go-gost/x/config"
)

type ProxyTransport struct {
	UseEncryption        bool   `json:"useEncryption,omitempty"`
	UseCompression       bool   `json:"useCompression,omitempty"`
	BandwidthLimit       string `json:"bandwidthLimit,omitempty"`
	BandwidthLimitMode   string `json:"bandwidthLimitMode,omitempty"`
	ProxyProtocolVersion string `json:"proxyProtocolVersion,omitempty"`
}

type HTTPProxyConfig struct {
	v1.HTTPProxyConfig
	Transport ProxyTransport
}

func (m HTTPProxyConfig) To() *v1.HTTPProxyConfig {
	quantity, _ := types.NewBandwidthQuantity(m.Transport.BandwidthLimit)
	m.HTTPProxyConfig.ProxyBaseConfig = v1.ProxyBaseConfig{
		Name:        m.ProxyBaseConfig.Name,
		Type:        m.ProxyBaseConfig.Type,
		Annotations: m.ProxyBaseConfig.Annotations,
		Transport: v1.ProxyTransport{
			UseEncryption:        m.Transport.UseEncryption,
			UseCompression:       m.Transport.UseCompression,
			BandwidthLimit:       quantity,
			BandwidthLimitMode:   m.Transport.BandwidthLimitMode,
			ProxyProtocolVersion: m.Transport.ProxyProtocolVersion,
		},
		Metadatas:    m.ProxyBaseConfig.Metadatas,
		LoadBalancer: m.ProxyBaseConfig.LoadBalancer,
		HealthCheck:  m.ProxyBaseConfig.HealthCheck,
		ProxyBackend: m.ProxyBaseConfig.ProxyBackend,
	}
	return &m.HTTPProxyConfig
}

type TCPProxyConfig struct {
	v1.TCPProxyConfig
	Transport ProxyTransport
}

func (m TCPProxyConfig) To() *v1.TCPProxyConfig {
	quantity, _ := types.NewBandwidthQuantity(m.Transport.BandwidthLimit)
	m.TCPProxyConfig.ProxyBaseConfig = v1.ProxyBaseConfig{
		Name:        m.ProxyBaseConfig.Name,
		Type:        m.ProxyBaseConfig.Type,
		Annotations: m.ProxyBaseConfig.Annotations,
		Transport: v1.ProxyTransport{
			UseEncryption:        m.Transport.UseEncryption,
			UseCompression:       m.Transport.UseCompression,
			BandwidthLimit:       quantity,
			BandwidthLimitMode:   m.Transport.BandwidthLimitMode,
			ProxyProtocolVersion: m.Transport.ProxyProtocolVersion,
		},
		Metadatas:    m.ProxyBaseConfig.Metadatas,
		LoadBalancer: m.ProxyBaseConfig.LoadBalancer,
		HealthCheck:  m.ProxyBaseConfig.HealthCheck,
		ProxyBackend: m.ProxyBaseConfig.ProxyBackend,
	}
	return &m.TCPProxyConfig
}

type UDPProxyConfig struct {
	v1.UDPProxyConfig
	Transport ProxyTransport
}

func (m UDPProxyConfig) To() *v1.UDPProxyConfig {
	quantity, _ := types.NewBandwidthQuantity(m.Transport.BandwidthLimit)
	m.UDPProxyConfig.ProxyBaseConfig = v1.ProxyBaseConfig{
		Name:        m.ProxyBaseConfig.Name,
		Type:        m.ProxyBaseConfig.Type,
		Annotations: m.ProxyBaseConfig.Annotations,
		Transport: v1.ProxyTransport{
			UseEncryption:        m.Transport.UseEncryption,
			UseCompression:       m.Transport.UseCompression,
			BandwidthLimit:       quantity,
			BandwidthLimitMode:   m.Transport.BandwidthLimitMode,
			ProxyProtocolVersion: m.Transport.ProxyProtocolVersion,
		},
		Metadatas:    m.ProxyBaseConfig.Metadatas,
		LoadBalancer: m.ProxyBaseConfig.LoadBalancer,
		HealthCheck:  m.ProxyBaseConfig.HealthCheck,
		ProxyBackend: m.ProxyBaseConfig.ProxyBackend,
	}
	return &m.UDPProxyConfig
}

type STCPProxyConfig struct {
	v1.STCPProxyConfig
	Transport ProxyTransport
}

func (m STCPProxyConfig) To() *v1.STCPProxyConfig {
	quantity, _ := types.NewBandwidthQuantity(m.Transport.BandwidthLimit)
	m.STCPProxyConfig.ProxyBaseConfig = v1.ProxyBaseConfig{
		Name:        m.ProxyBaseConfig.Name,
		Type:        m.ProxyBaseConfig.Type,
		Annotations: m.ProxyBaseConfig.Annotations,
		Transport: v1.ProxyTransport{
			UseEncryption:        m.Transport.UseEncryption,
			UseCompression:       m.Transport.UseCompression,
			BandwidthLimit:       quantity,
			BandwidthLimitMode:   m.Transport.BandwidthLimitMode,
			ProxyProtocolVersion: m.Transport.ProxyProtocolVersion,
		},
		Metadatas:    m.ProxyBaseConfig.Metadatas,
		LoadBalancer: m.ProxyBaseConfig.LoadBalancer,
		HealthCheck:  m.ProxyBaseConfig.HealthCheck,
		ProxyBackend: m.ProxyBaseConfig.ProxyBackend,
	}
	return &m.STCPProxyConfig
}

type SUDPProxyConfig struct {
	v1.SUDPProxyConfig
	Transport ProxyTransport
}

func (m SUDPProxyConfig) To() *v1.SUDPProxyConfig {
	quantity, _ := types.NewBandwidthQuantity(m.Transport.BandwidthLimit)
	m.SUDPProxyConfig.ProxyBaseConfig = v1.ProxyBaseConfig{
		Name:        m.ProxyBaseConfig.Name,
		Type:        m.ProxyBaseConfig.Type,
		Annotations: m.ProxyBaseConfig.Annotations,
		Transport: v1.ProxyTransport{
			UseEncryption:        m.Transport.UseEncryption,
			UseCompression:       m.Transport.UseCompression,
			BandwidthLimit:       quantity,
			BandwidthLimitMode:   m.Transport.BandwidthLimitMode,
			ProxyProtocolVersion: m.Transport.ProxyProtocolVersion,
		},
		Metadatas:    m.ProxyBaseConfig.Metadatas,
		LoadBalancer: m.ProxyBaseConfig.LoadBalancer,
		HealthCheck:  m.ProxyBaseConfig.HealthCheck,
		ProxyBackend: m.ProxyBaseConfig.ProxyBackend,
	}
	return &m.SUDPProxyConfig
}

type XTCPProxyConfig struct {
	v1.XTCPProxyConfig
	Transport ProxyTransport
}

func (m XTCPProxyConfig) To() *v1.XTCPProxyConfig {
	quantity, _ := types.NewBandwidthQuantity(m.Transport.BandwidthLimit)
	m.XTCPProxyConfig.ProxyBaseConfig = v1.ProxyBaseConfig{
		Name:        m.ProxyBaseConfig.Name,
		Type:        m.ProxyBaseConfig.Type,
		Annotations: m.ProxyBaseConfig.Annotations,
		Transport: v1.ProxyTransport{
			UseEncryption:        m.Transport.UseEncryption,
			UseCompression:       m.Transport.UseCompression,
			BandwidthLimit:       quantity,
			BandwidthLimitMode:   m.Transport.BandwidthLimitMode,
			ProxyProtocolVersion: m.Transport.ProxyProtocolVersion,
		},
		Metadatas:    m.ProxyBaseConfig.Metadatas,
		LoadBalancer: m.ProxyBaseConfig.LoadBalancer,
		HealthCheck:  m.ProxyBaseConfig.HealthCheck,
		ProxyBackend: m.ProxyBaseConfig.ProxyBackend,
	}
	return &m.XTCPProxyConfig
}

type HostConfig struct {
	Key     string // 唯一标识
	BaseCfg v1.ClientCommonConfig
	Http    HTTPProxyConfig
}

type ForwardConfig struct {
	Key     string // 唯一标识
	BaseCfg v1.ClientCommonConfig
	TCP     TCPProxyConfig
	UDP     UDPProxyConfig
}

type TunnelConfig struct {
	Key     string // 唯一标识
	BaseCfg v1.ClientCommonConfig
	STCP    STCPProxyConfig
	SUDP    SUDPProxyConfig
}

type ProxyConfig struct {
	Key      string // 唯一标识
	BaseCfg  v1.ClientCommonConfig
	Name     string
	Port     int
	AuthUser string
	AuthPwd  string
	Metadata map[string]string
	Limiter  string
}

type P2PConfig struct {
	Key     string // 唯一标识
	BaseCfg v1.ClientCommonConfig
	XTCP    XTCPProxyConfig
	STCP    STCPProxyConfig
}

type ServerConfig struct {
	Key string // 唯一标识
	v1.ServerConfig
}

type ServerDomain struct {
	Domain     string
	Target     string
	Cert       string
	Key        string
	ForceHttps int
}

type ServerDomainCacheData struct {
	Domain            string
	Active            bool
	CacheRules        []ServerDomainCacheRuleConfig
	CacheExcludeRules []ServerDomainCacheRuleConfig
}

type ServerDomainCacheRuleConfig struct {
	Method string
	Type   string // 规则类型 prefix:前缀匹配 suffix:后缀匹配 include:包含 full:全部匹配
	Rule   string
}

type TunMasterConfig struct {
	Key            string
	Service        config.ServiceConfig
	ForwardService config.ServiceConfig
	Auther         config.AutherConfig
}

type TunNodeConfig struct {
	Key     string
	Service config.ServiceConfig
	Chain   config.ChainConfig
}

type CustomCfgConfig struct {
	Key     string
	Type    string
	Content string
}
