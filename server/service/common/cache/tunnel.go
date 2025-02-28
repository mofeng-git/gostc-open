package cache

import (
	"github.com/patrickmn/go-cache"
	"server/global"
)

const (
	tunnel_info_key = "tunnel:info:"
)

type TunnelInfo struct {
	Code        string
	Type        int
	ClientCode  string
	UserCode    string
	NodeCode    string
	ChargingTye int
	ExpAt       int64
	Limiter     int
}

func SetTunnelInfo(info TunnelInfo) {
	global.Cache.SetStruct(tunnel_info_key+info.Code, info, cache.NoExpiration)
}

func GetTunnelInfo(code string) (result TunnelInfo) {
	_ = global.Cache.GetStruct(tunnel_info_key+code, &result)
	return result
}

func DelTunnelInfo(code string) {
	global.Cache.Del(tunnel_info_key + code)
}
