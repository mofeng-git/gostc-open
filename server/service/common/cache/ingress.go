package cache

import (
	"github.com/patrickmn/go-cache"
	"server/global"
)

const (
	gost_plugins_ingress_key = "gost_plugins:ingress:"
)

func SetIngress(host, tunnelCode string) {
	global.Cache.SetString(gost_plugins_ingress_key+host, tunnelCode, cache.NoExpiration)
}

func GetIngress(host string) string {
	return global.Cache.GetString(gost_plugins_ingress_key + host)
}

func DelIngress(host string) {
	global.Cache.Del(gost_plugins_ingress_key + host)
}
