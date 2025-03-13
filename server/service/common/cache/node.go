package cache

import (
	"github.com/patrickmn/go-cache"
	"server/global"
	"sync"
	"time"
)

const (
	node_online_key   = "node:online:"
	node_version_key  = "node:version:"
	node_port_use_key = "node:port:"
)

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
