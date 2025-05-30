package cache

import (
	"github.com/patrickmn/go-cache"
	"server/global"
	"time"
)

const (
	auth_bind_email_key                = "auth:bind:email:"
	auth_bind_email_last_time_key      = "auth:bind:email:last:time:"
	auth_reset_pwd_email_key           = "auth:reset:pwd:email:"
	auth_reset_pwd_email_last_time_key = "auth:reset:pwd:email:last:time:"
)

// 记录用户上一次发送重置密码验证码时间
func SetResetPwdEmailLastTime(code string) {
	global.Cache.SetString(auth_reset_pwd_email_last_time_key+code, time.Now().Format(time.DateTime), cache.NoExpiration)
}

// 获取用户上一次发送重置密码验证码时间
func GetResetPwdEmailLastTime(code string) time.Time {
	value := global.Cache.GetString(auth_reset_pwd_email_last_time_key + code)
	if value == "" {
		return time.Time{}
	}
	location, _ := time.ParseInLocation(time.DateTime, value, time.Local)
	return location
}

func SetResetPwdEmailCode(key string, value string, duration time.Duration) {
	global.Cache.SetString(auth_reset_pwd_email_key+key, value, duration)
}

func GetResetPwdEmailCode(key string, remove bool) string {
	defer func() {
		if remove {
			global.Cache.Del(auth_reset_pwd_email_key + key)
		}
	}()
	return global.Cache.GetString(auth_reset_pwd_email_key + key)
}

// 记录用户上一次发送绑定邮箱验证码时间
func SetBindEmailLastTime(code string) {
	global.Cache.SetString(auth_bind_email_last_time_key+code, time.Now().Format(time.DateTime), cache.NoExpiration)
}

// 获取用户上一次发送绑定邮箱验证码时间
func GetBindEmailLastTime(code string) time.Time {
	value := global.Cache.GetString(auth_bind_email_last_time_key + code)
	if value == "" {
		return time.Time{}
	}
	location, _ := time.ParseInLocation(time.DateTime, value, time.Local)
	return location
}

func SetBindEmailCode(key string, value string, duration time.Duration) {
	global.Cache.SetString(auth_bind_email_key+key, value, duration)
}

func GetBindEmailCode(key string, remove bool) string {
	defer func() {
		if remove {
			global.Cache.Del(auth_bind_email_key + key)
		}
	}()
	return global.Cache.GetString(auth_bind_email_key + key)
}
