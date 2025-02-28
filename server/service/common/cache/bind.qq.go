package cache

import (
	"server/global"
	"time"
)

const (
	auth_bind_qq_key = "auth:bind:qq:"
)

func SetBindQQCode(key string, value string, duration time.Duration) {
	global.Cache.SetString(auth_bind_qq_key+key, value, duration)
}

func GetBindQQCode(key string, remove bool) string {
	defer func() {
		if remove {
			global.Cache.Del(auth_bind_qq_key + key)
		}
	}()
	return global.Cache.GetString(auth_bind_qq_key + key)
}
