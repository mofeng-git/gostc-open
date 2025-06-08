package cache

import (
	"server/global"
	"time"
)

const (
	auth_login_otp_key = "auth:login:otp:"
)

func SetLoginOtp(key string, value string, duration time.Duration) {
	global.Cache.SetString(auth_login_otp_key+key, value, duration)
}

func GetLoginOtp(key string, remove bool) string {
	defer func() {
		if remove {
			global.Cache.Del(auth_login_otp_key + key)
		}
	}()
	return global.Cache.GetString(auth_login_otp_key + key)
}

func DelLoginOtp(key string) {
	global.Cache.Del(auth_login_otp_key + key)
}
