package cache

import (
	"github.com/patrickmn/go-cache"
	"server/global"
	"sync"
	"time"
)

const (
	node_online_key        = "node:online:"
	node_version_key       = "node:version:"
	node_port_use_key      = "node:port:"
	node_custom_domain_key = "node:custom:domain:"
	node_cache_key         = "node:cache:"
	node_info_key          = "node:info:"
)

type NodeInfo struct {
	Code            string
	LimitResetIndex int
	LimitTotal      int
	LimitKind       int
}

func SetNodeInfo(data NodeInfo) {
	global.Cache.SetStruct(node_info_key+data.Code, data, cache.NoExpiration)
}

func GetNodeInfo(code string) (result NodeInfo) {
	_ = global.Cache.GetStruct(node_info_key+code, &result)
	return result
}

func DelNodeInfo(code string) {
	global.Cache.Del(node_info_key + code)
}

func SetNodeOnline(code string, online bool, duration time.Duration) {
	global.Cache.SetString(node_online_key+code, func(online bool) string {
		if online {
			return "1"
		}
		return "2"
	}(online), duration)
}

func GetNodeOnline(code string) bool {
	return global.Cache.GetString(node_online_key+code) == "1"
}

func SetNodeCustomDomain(code string, enable bool) {
	global.Cache.SetString(node_custom_domain_key+code, func(enable bool) string {
		if enable {
			return "1"
		}
		return "2"
	}(enable), cache.NoExpiration)
}

func GetNodeCustomDomain(code string) bool {
	return global.Cache.GetString(node_custom_domain_key+code) == "1"
}

func SetNodeCache(code string, enable bool) {
	global.Cache.SetString(node_cache_key+code, func(enable bool) string {
		if enable {
			return "1"
		}
		return "2"
	}(enable), cache.NoExpiration)
}

func GetNodeCache(code string) bool {
	return global.Cache.GetString(node_cache_key+code) == "1"
}

func SetNodeVersion(code string, version string) {
	global.Cache.SetString(node_version_key+code, version, cache.NoExpiration)
}

func GetNodeVersion(code string) string {
	return global.Cache.GetString(node_version_key + code)
}

var nodePortUseLock = &sync.RWMutex{}

func SetNodePortUse(code string, port string, use bool, duration time.Duration) {
	nodePortUseLock.Lock()
	defer nodePortUseLock.Unlock()
	var data = make(map[string]bool)
	_ = global.Cache.GetStruct(node_port_use_key+code, &data)
	data[port] = use
	global.Cache.SetStruct(node_port_use_key+code, data, duration)
}

func GetNodePortUse(code string, port string) (bool, bool) {
	nodePortUseLock.RLock()
	defer nodePortUseLock.RUnlock()
	var data = make(map[string]bool)
	_ = global.Cache.GetStruct(node_port_use_key+code, &data)
	use, exist := data[port]
	return use, exist
}

func DelNodePortUse(code string) {
	global.Cache.Del(node_port_use_key + code)
}
