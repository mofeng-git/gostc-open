package service

import (
	"server/model"
	"server/service/common/cache"
	"time"
)

type LimiterReq struct {
	Scope   string `json:"scope"`
	Service string `json:"service"`
	Network string `json:"network"`
	Addr    string `json:"addr"`
	Client  string `json:"client"`
	Src     string `json:"src"`
}
type LimiterResp struct {
	In  int `json:"in"`
	Out int `json:"out"`
}

func (service *service) Limiter(req LimiterReq) LimiterResp {
	if req.Scope != "client" {
		return LimiterResp{In: 1024, Out: 1024}
	}
	var speed int
	result := cache.GetTunnelInfo(req.Client)
	if result.Code == "" || (result.ExpAt < time.Now().Unix() && result.ChargingTye == model.GOST_CONFIG_CHARGING_CUCLE_DAY) {
		speed = 1024
	} else {
		speed = result.Limiter * 128 * 1024
	}
	return LimiterResp{
		In:  speed,
		Out: speed,
	}
}
