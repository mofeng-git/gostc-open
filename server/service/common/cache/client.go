package cache

import (
	"github.com/patrickmn/go-cache"
	"server/global"
	"time"
)

const (
	client_online_key    = "client:online:"
	client_ip_key        = "client:ip:"
	client_last_time_key = "client:last_time:"
	client_version_key   = "client:version:"
)

func SetClientOnline(code string, online bool, duration time.Duration) {
	global.Cache.SetString(client_online_key+code, func(online bool) string {
		if online {
			return "1"
		}
		return "2"
	}(online), duration)
}

func GetClientOnline(code string) bool {
	return global.Cache.GetString(client_online_key+code) == "1"
}

func SetClientIp(code string, ip string) {
	global.Cache.SetString(client_ip_key+code, ip, cache.NoExpiration)
}

func GetClientIp(code string) string {
	return global.Cache.GetString(client_ip_key + code)
}

func SetClientVersion(code string, version string) {
	global.Cache.SetString(client_version_key+code, version, cache.NoExpiration)
}

func GetClientVersion(code string) string {
	return global.Cache.GetString(client_version_key + code)
}

func SetClientLastTime(code string) {
	global.Cache.SetString(client_last_time_key+code, time.Now().Format(time.DateTime), cache.NoExpiration)
}

func GetClientLastTime(code string) string {
	return global.Cache.GetString(client_last_time_key + code)
}
