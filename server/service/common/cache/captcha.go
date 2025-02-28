package cache

import (
	"server/global"
	"time"
)

const (
	auth_captcha_key = "auth:captcha:"
)

func SetCaptcha(key string, value string, duration time.Duration) {
	global.Cache.SetString(auth_captcha_key+key, value, duration)
}

func ValidCaptcha(key string, target string, remove bool) bool {
	defer func() {
		if remove {
			global.Cache.Del(auth_captcha_key + key)
		}
	}()
	return global.Cache.GetString(auth_captcha_key+key) == target
}
