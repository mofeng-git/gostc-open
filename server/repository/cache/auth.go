package cache

import (
	"github.com/patrickmn/go-cache"
	"server/global"
)

const (
	gost_auth_key = "gost:auth:"
)

func SetGostAuth(user, pwd, code string) {
	global.Cache.SetString(gost_auth_key+user+pwd, code, cache.NoExpiration)
}

func GetGostAuth(user, pwd string) string {
	return global.Cache.GetString(gost_auth_key + user + pwd)
}

func DelGostAuth(user, pwd string) {
	global.Cache.Del(gost_auth_key + user + pwd)
}
