package model

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	GOST_NODE_LIMIT_KIND_ALL    = 0 // 全部流量
	GOST_NODE_LIMIT_KIND_INPUT  = 1 // 上行流量
	GOST_NODE_LIMIT_KIND_OUTPUT = 2 // 下行流量
)

type GostNode struct {
	Base
	IndexValue         int              `gorm:"column:index_value;index;default:1000;comment:排序，升序"`
	Key                string           `gorm:"column:key;size:100;uniqueIndex"`
	Name               string           `gorm:"column:name;index;comment:名称"`
	Remark             string           `gorm:"column:remark;comment:节点介绍"`
	Web                int              `gorm:"column:web;size:1;default:2;comment:WEB功能"`
	Tunnel             int              `gorm:"column:tunnel;size:1;default:2;comment:私有隧道功能"`
	Forward            int              `gorm:"column:forward;size:1;default:2;comment:端口转发功能"`
	Proxy              int              `gorm:"column:proxy;size:1;default:2;comment:代理隧道功能"`
	P2P                int              `gorm:"column:p2p;size:1;default:2;comment:P2P隧道功能"`
	Domain             string           `gorm:"column:domain;comment:基础域名"`
	DenyDomainPrefix   string           `gorm:"column:deny_domain_prefix;comment:不允许的域名前缀"`
	AllowDomainMatcher int              `gorm:"column:allow_domain_matcher;comment:自定义域名是否允许泛域名"`
	UrlTpl             string           `gorm:"column:url_tpl;comment:URL模板"`
	Protocol           string           `gorm:"column:protocol;comment:协议"`
	Address            string           `gorm:"column:address;comment:服务地址"`
	HttpPort           string           `gorm:"column:http_port;comment:HTTP流量端口"`
	ReplaceAddress     string           `gorm:"column:replace_address;comment:替换地址"`
	ForwardPorts       string           `gorm:"column:forward_ports;comment:转发端口组"`
	P2PDisableForward  int              `gorm:"column:p2p_disable_forward;size:1;default:0;comment:是否禁用中继"`
	Rules              string           `gorm:"column:rules;comment:规则限制"`
	Tags               string           `gorm:"column:tags;comment:标签"`
	Configs            []GostNodeConfig `gorm:"foreignKey:NodeCode;references:Code"`

	LimitResetIndex int `gorm:"column:limit_reset_index;comment:重置日期索引"`
	LimitTotal      int `gorm:"column:limit_total;comment:预警流量(GB)"`
	LimitKind       int `gorm:"column:limit_kind;size:1;default:0;comment:限制方式"`
}

// 生成更新指纹，对比指纹判断是否需要标记更新
func (n GostNode) GenerateFingerprint() string {
	return fmt.Sprintf("%s|%s|%s|%s|%s|%d", n.Domain, n.Protocol, n.Address, n.HttpPort, n.ReplaceAddress, n.AllowDomainMatcher)
}

func (n GostNode) GetAddress() (host string, port int) {
	var address = n.Address
	if n.ReplaceAddress != "" {
		address = n.ReplaceAddress
	}
	splitHost, splitPort, _ := net.SplitHostPort(address)
	port, _ = strconv.Atoi(splitPort)
	return splitHost, port
}

func (n GostNode) GetOriginAddress() (host string, port int) {
	var address = n.Address
	splitHost, splitPort, _ := net.SplitHostPort(address)
	port, _ = strconv.Atoi(splitPort)
	return splitHost, port
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

func (n GostNode) GetDomainHost(domainPrefix string, customDomain string, enableCustom bool) string {
	if enableCustom && customDomain != "" {
		return customDomain
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
