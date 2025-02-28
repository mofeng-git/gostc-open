package cache

import (
	"server/global"
	"time"
)

const (
	auth_bind_otp_key = "auth:bind:otp:"
)

func SetBindOtp(key string, value string, duration time.Duration) {
	global.Cache.SetString(auth_bind_otp_key+key, value, duration)
}

func GetBindOtp(key string, remove bool) string {
	defer func() {
		if remove {
			global.Cache.Del(auth_bind_otp_key + key)
		}
	}()
	return global.Cache.GetString(auth_bind_otp_key + key)
}
