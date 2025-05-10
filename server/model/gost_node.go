package model

import (
	"encoding/json"
	"fmt"
	"github.com/go-gost/x/config"
	v1 "server/pkg/p2p_cfg/v1"
	"strconv"
	"strings"
)

type GostNode struct {
	Base
	IndexValue            int              `gorm:"column:index_value;index;default:1000;comment:排序，升序"`
	Key                   string           `gorm:"column:key;size:100;uniqueIndex"`
	Name                  string           `gorm:"column:name;index;comment:名称"`
	Remark                string           `gorm:"column:remark;comment:节点介绍"`
	Web                   int              `gorm:"column:web;size:1;default:2;comment:WEB功能"`
	Tunnel                int              `gorm:"column:tunnel;size:1;default:2;comment:私有隧道功能"`
	Forward               int              `gorm:"column:forward;size:1;default:2;comment:端口转发功能"`
	Proxy                 int              `gorm:"column:proxy;size:1;default:2;comment:代理隧道功能"`
	P2P                   int              `gorm:"column:p2p;size:1;default:2;comment:P2P隧道功能"`
	Domain                string           `gorm:"column:domain;comment:基础域名"`
	DenyDomainPrefix      string           `gorm:"column:deny_domain_prefix;comment:不允许的域名前缀"`
	UrlTpl                string           `gorm:"column:url_tpl;comment:URL模板"`
	Address               string           `gorm:"column:address;comment:服务地址"`
	Protocol              string           `gorm:"column:protocol;comment:协议"`
	TunnelConnPort        string           `gorm:"column:tunnel_conn_port;comment:隧道连接端口"`
	TunnelInPort          string           `gorm:"column:tunnel_in_port;comment:隧道访问端口"`
	TunnelMetadata        string           `gorm:"column:tunnel_metadata;comment:其他信息"`
	TunnelReplaceAddress  string           `gorm:"column:tunnel_replace_address;comment:替换地址"`
	ForwardConnPort       string           `gorm:"column:forward_conn_port;comment:转发连接端口"`
	ForwardPorts          string           `gorm:"column:forward_ports;comment:转发端口组"`
	ForwardMetadata       string           `gorm:"column:forward_metadata;comment:其他信息"`
	ForwardReplaceAddress string           `gorm:"column:forward_replace_address;comment:替换地址"`
	P2PPort               string           `gorm:"column:p2p_port;comment:p2p连接端口"`
	Rules                 string           `gorm:"column:rules;comment:规则限制"`
	Tags                  string           `gorm:"column:tags;comment:标签"`
	Configs               []GostNodeConfig `gorm:"foreignKey:NodeCode;references:Code"`
}

func (n GostNode) GetDomainFull(domainPrefix string, customDomain string, enableCustom bool) string {
	if enableCustom && customDomain != "" {
		if n.UrlTpl == "" {
			return customDomain
		}
		return strings.ReplaceAll(n.UrlTpl, "{{DOMAIN}}", customDomain)
	}
	if n.UrlTpl != "" {
		return strings.ReplaceAll(n.UrlTpl, "{{DOMAIN}}", domainPrefix+"."+n.Domain)
	}
	return domainPrefix + "." + n.Domain
}

func (n GostNode) GetRules() (result []string) {
	for _, rule := range strings.Split(n.Rules, ",") {
		if rule == "" {
			continue
		}
		result = append(result, rule)
	}
	if len(result) == 0 {
		result = append(result, "")
	}
	return result
}

func (n GostNode) GetTags() (result []string) {
	for _, tag := range strings.Split(n.Tags, ",") {
		if tag == "" {
			continue
		}
		result = append(result, tag)
	}
	if len(result) == 0 {
		result = append(result, "暂无标签")
	}
	return
}

func (n GostNode) CheckDomainPrefix(prefix string) bool {
	for _, allow := range strings.Split(n.DenyDomainPrefix, "\n") {
		if allow == "" {
			continue
		}
		if prefix == allow {
			return false
		}
	}
	return true
}

func (n GostNode) GetPorts(excludePort []string) (result map[string]bool) {
	result = make(map[string]bool)
	var excludePortMap = make(map[string]bool)
	for _, port := range excludePort {
		excludePortMap[port] = true
	}
	for _, v1 := range strings.Split(strings.ReplaceAll(n.ForwardPorts, " ", ""), ",") {
		if v1 == "" {
			continue
		}
		if _, err := strconv.Atoi(v1); err == nil {
			result[v1] = true
			continue
		}
		portGroup := strings.Split(v1, "-")
		if len(portGroup) != 2 {
			continue
		}
		start, err := strconv.Atoi(portGroup[0])
		if err != nil {
			continue
		}
		end, err := strconv.Atoi(portGroup[1])
		if err != nil {
			continue
		}
		if start >= end {
			continue
		}
		for {
			if start > end {
				break
			}
			result[strconv.Itoa(start)] = true
			start++
		}
	}
	for k, _ := range excludePortMap {
		delete(result, k)
	}
	return result
}

func (n GostNode) GenerateTunnelAndHostServiceConfig(limiter, auther, ingress, obs string) (config.ServiceConfig, bool) {
	if n.Tunnel != 1 {
		return config.ServiceConfig{}, false
	}
	var metadata = make(map[string]any)
	_ = json.Unmarshal([]byte(n.TunnelMetadata), &metadata)
	if n.Web == 1 {
		metadata["entrypoint"] = ":" + n.TunnelInPort
	}
	metadata["ingress"] = ingress
	metadata["sniffing"] = true
	metadata["limiter.scope"] = "service"
	metadata["observer.period"] = "60s"
	metadata["observer.resetTraffic"] = true
	return config.ServiceConfig{
		Name: "tunnel_" + n.Code,
		Addr: ":" + n.TunnelConnPort,
		Handler: &config.HandlerConfig{
			Type:     "tunnel",
			Metadata: metadata,
			Auther:   auther,
			Observer: obs,
			Limiter:  limiter,
		},
		Listener: &config.ListenerConfig{
			Type: n.Protocol,
		},
	}, true
}

func (n GostNode) GenerateForwardServiceConfig(limiter, auther, obs string) (config.ServiceConfig, bool) {
	if n.Forward != 1 {
		return config.ServiceConfig{}, false
	}
	var metadata = make(map[string]any)
	_ = json.Unmarshal([]byte(n.TunnelMetadata), &metadata)
	metadata["bind"] = true
	//metadata["nodelay"] = true
	metadata["limiter.scope"] = "service"
	metadata["observer.period"] = "60s"
	metadata["observer.resetTraffic"] = true
	return config.ServiceConfig{
		Name: "forward_" + n.Code,
		Addr: ":" + n.ForwardConnPort,
		Handler: &config.HandlerConfig{
			Type:     "relay",
			Metadata: metadata,
			Auther:   auther,
			Limiter:  limiter,
			Observer: obs,
		},
		Listener: &config.ListenerConfig{
			Type: n.Protocol,
		},
	}, true
}

func (n GostNode) GenerateIngress(hosts []*GostClientHost, tunnels []*GostClientTunnel, enableCustom bool) config.IngressConfig {
	var rules []*config.IngressRuleConfig
	for _, host := range hosts {
		hostname := host.DomainPrefix + "." + n.Domain
		if host.CustomDomain != "" && enableCustom {
			hostname = host.CustomDomain
		}
		rules = append(rules, &config.IngressRuleConfig{
			Hostname: hostname,
			Endpoint: host.Code,
		})
	}
	for _, tunnel := range tunnels {
		rules = append(rules, &config.IngressRuleConfig{
			Hostname: tunnel.Code,
			Endpoint: "$" + tunnel.Code,
		})
	}
	return config.IngressConfig{
		Name:  "ingress_" + n.Code,
		Rules: rules,
	}
	//return config.IngressConfig{
	//	Name: "ingress_" + n.Code,
	//	Plugin: &config.PluginConfig{
	//		Type: "http",
	//		Addr: fmt.Sprintf("%s/api/v1/public/gost/ingress", host),
	//	},
	//}
}

func (n GostNode) GenerateLimiter(host string) config.LimiterConfig {
	return config.LimiterConfig{
		Name: "limiter_" + n.Code,
		Plugin: &config.PluginConfig{
			Type: "http",
			Addr: fmt.Sprintf("%s/api/v1/public/gost/limiter", host),
		},
	}
}

func (n GostNode) GenerateAuther(host string) config.AutherConfig {
	return config.AutherConfig{
		Name: "auther_" + n.Code,
		Plugin: &config.PluginConfig{
			Type: "http",
			Addr: fmt.Sprintf("%s/api/v1/public/gost/auther", host),
		},
	}
}

func (n GostNode) GenerateObs(host string) config.ObserverConfig {
	return config.ObserverConfig{
		Name: "obs_" + n.Code,
		Plugin: &config.PluginConfig{
			Type: "http",
			Addr: fmt.Sprintf("%s/api/v1/public/gost/obs", host),
		},
	}
}

func (n GostNode) GenerateNodePortCheck(host string, port string) map[string]string {
	return map[string]string{
		"callback": fmt.Sprintf("%s/api/v1/public/gost/node/port", host),
		"code":     n.Code,
		"port":     port,
	}
}

func (n GostNode) GenerateP2PServiceConfig(host string) v1.ServerConfig {
	port, _ := strconv.Atoi(n.P2PPort)
	return v1.ServerConfig{
		BindPort: port,
		HTTPPlugins: []v1.HTTPPluginOptions{
			//{
			//	Name:      "login-plugin",
			//	Addr:      host,
			//	Path:      "/api/v1/public/p2p/login",
			//	Ops:       []string{"Login"},
			//	TLSVerify: true,
			//},
			{
				Name:      "new-plugin",
				Addr:      host,
				Path:      "/api/v1/public/p2p/new",
				Ops:       []string{"NewProxy"},
				TLSVerify: true,
			},
		},
	}
}
