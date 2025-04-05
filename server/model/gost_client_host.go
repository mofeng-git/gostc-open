package model

import (
	"encoding/json"
	"github.com/go-gost/x/config"
	"strings"
)

type GostClientHost struct {
	Base
	Name         string     `gorm:"column:name;index;comment:名称"`
	TargetIp     string     `gorm:"column:target_ip;index;comment:内网IP"`
	TargetPort   string     `gorm:"column:target_port;index;comment:内网端口"`
	DomainPrefix string     `gorm:"column:domain_prefix;index;comment:域名前缀"`
	CustomDomain string     `gorm:"column:custom_domain;index;comment:自定义域名"`
	CustomCert   string     `gorm:"column:custom_cert;comment:自定义证书"`
	CustomKey    string     `gorm:"column:custom_key;comment:自定义证书"`
	NodeCode     string     `gorm:"column:node_code;index;comment:节点编号"`
	Node         GostNode   `gorm:"foreignKey:NodeCode;references:Code"`
	ClientCode   string     `gorm:"column:client_code;index;comment:客户端编号"`
	Client       GostClient `gorm:"foreignKey:ClientCode;references:Code"`
	UserCode     string     `gorm:"column:user_code;index;comment:用户编号"`
	User         SystemUser `gorm:"foreignKey:UserCode;references:Code"`
	Enable       int        `gorm:"column:enable;size:1;default:1;comment:启用状态"`
	Status       int        `gorm:"column:status;size:1;default:1;comment:状态"`
	GostClientAdmission
	GostClientConfig
}

func (host GostClientHost) GenerateSvcConfig(chain, admissionWhite, admissionBlack string) (config.ServiceConfig, bool) {
	if host.Node.Web != 1 {
		return config.ServiceConfig{}, false
	}
	var target = host.TargetIp + ":" + host.TargetPort
	var admissions []string
	if host.WhiteEnable == 1 {
		admissions = append(admissions, admissionWhite)
	}
	if host.BlackEnable == 1 {
		admissions = append(admissions, admissionBlack)
	}
	return config.ServiceConfig{
		Name:       host.Code,
		Addr:       ":0",
		Observer:   "",
		Admissions: admissions,
		Recorders:  nil,
		Handler:    &config.HandlerConfig{Type: "rtcp", Metadata: map[string]any{"sniffing": true}},
		Listener:   &config.ListenerConfig{Type: "rtcp", Chain: chain},
		Forwarder: &config.ForwarderConfig{
			Nodes: []*config.ForwardNodeConfig{
				{
					Name: target,
					Addr: target,
					HTTP: &config.HTTPNodeConfig{
						Host:           host.DomainPrefix + "." + host.Node.Domain,
						ResponseHeader: nil,
						Auth:           nil,
					},
				},
				{
					Name: target,
					Addr: target,
					HTTP: &config.HTTPNodeConfig{
						Host: host.DomainPrefix + "." + host.Node.Domain,
						ResponseHeader: map[string]string{
							"Cache-Control": "private, max-age=86400, stale-while-revalidate=604800",
						},
						Auth: nil,
					},
					Matcher: &config.NodeMatcherConfig{
						Rule: "!Header(`key`) && PathRegexp(`\\.(css|js|jpeg|jpg|png|gif|bmp|webp|svg|ttf|woff|woff2)$`)",
					},
				},
			},
		},
	}, true
}

func (host GostClientHost) GenerateChainConfig(auth GostAuth) config.ChainConfig {
	var metadata = make(map[string]any)
	_ = json.Unmarshal([]byte(host.Node.TunnelMetadata), &metadata)
	metadata["tunnel.id"] = host.Code

	var protocol, address string
	protocol = host.Node.Protocol
	address = host.Node.Address + ":" + host.Node.TunnelConnPort
	replaceAddress := strings.Split(host.Node.TunnelReplaceAddress, "://")
	if len(replaceAddress) == 2 {
		protocol = replaceAddress[0]
		address = replaceAddress[1]
	}

	return config.ChainConfig{
		Name: host.Code,
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

func (host GostClientHost) GenerateWhiteAdmission() config.AdmissionConfig {
	if host.WhiteEnable == 2 {
		return config.AdmissionConfig{}
	}
	return config.AdmissionConfig{
		Name:      "admission_white_" + host.Code,
		Whitelist: true,
		Matchers:  host.GetWhiteList(),
	}
}

func (host GostClientHost) GenerateBlackAdmission() config.AdmissionConfig {
	if host.BlackEnable == 2 {
		return config.AdmissionConfig{}
	}
	return config.AdmissionConfig{
		Name:      "admission_black_" + host.Code,
		Whitelist: false,
		Matchers:  host.GetBlackList(),
	}
}
