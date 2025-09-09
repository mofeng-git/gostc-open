package model

import "fmt"

type GostClientHost struct {
	Base
	Name                string     `gorm:"column:name;index;comment:名称"`
	TargetIp            string     `gorm:"column:target_ip;index;comment:内网IP"`
	TargetPort          string     `gorm:"column:target_port;index;comment:内网端口"`
	TargetHttps         int        `gorm:"column:target_https;comment:是否为Https服务"`
	DomainPrefix        string     `gorm:"column:domain_prefix;index;comment:域名前缀"`
	CustomDomain        string     `gorm:"column:custom_domain;index;comment:自定义域名"`
	CustomCert          string     `gorm:"column:custom_cert;comment:自定义证书"`
	CustomKey           string     `gorm:"column:custom_key;comment:自定义证书"`
	CustomForceHttps    int        `gorm:"column:custom_force_https;size:1;comment:强制HTTPS"`
	CustomDomainMatcher int        `gorm:"column:custom_domain_matcher;size:1;comment:是否为泛域名"`
	NodeCode            string     `gorm:"column:node_code;index;comment:节点编号"`
	Node                GostNode   `gorm:"foreignKey:NodeCode;references:Code"`
	ClientCode          string     `gorm:"column:client_code;index;comment:客户端编号"`
	Client              GostClient `gorm:"foreignKey:ClientCode;references:Code"`
	UserCode            string     `gorm:"column:user_code;index;comment:用户编号"`
	User                SystemUser `gorm:"foreignKey:UserCode;references:Code"`
	Enable              int        `gorm:"column:enable;size:1;default:1;comment:启用状态"`
	Status              int        `gorm:"column:status;size:1;default:1;comment:状态"`
	GostClientAdmission
	GostClientConfig
}

func (host *GostClientHost) GetCustomDomain() string {
	if host.CustomDomain == "" {
		return ""
	} else {
		if host.CustomDomainMatcher == 1 && host.Node.AllowDomainMatcher == 1 {
			return fmt.Sprintf("*.%s", host.CustomDomain)
		} else {
			return host.CustomDomain
		}
	}
}
