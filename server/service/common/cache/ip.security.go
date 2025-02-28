package cache

import (
	"github.com/patrickmn/go-cache"
	"server/global"
)

const (
	ip_security_key = "ip:security:"
)

func SetIpSecurity(ip string, security bool) {
	global.Cache.SetString(ip_security_key+ip, func() string {
		if security {
			return "1"
		}
		return "2"
	}(), cache.NoExpiration)
}

func GetIpSecurity(ip string) bool {
	return global.Cache.GetString(ip_security_key+ip) != "2"
}
