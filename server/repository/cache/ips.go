package cache

import (
	"errors"
	"net"
	"server/global"
	"time"
)

const (
	domain_ips_key = "domain:ips:"
)

func GetDomainIp(domain string) (string, error) {
	ip := global.Cache.GetString(domain_ips_key + domain)
	if ip == "" {
		ips, err := net.LookupHost(domain)
		if err != nil {
			return "", err
		}
		if len(ips) == 0 {
			return "", errors.New("该域名没有解析到指定IP，无法绑定")
		}
		if len(ips) > 1 {
			return "", errors.New("该域名被解析到了多个IP，无法绑定")
		}
		ip = ips[0]
		global.Cache.SetString(domain_ips_key+domain, ip, time.Minute*3)
	}
	return ip, nil
}
