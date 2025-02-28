package cache

import (
	"github.com/patrickmn/go-cache"
	"server/global"
	"time"
)

const (
	node_online_key  = "node:online:"
	node_version_key = "node:version:"
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
