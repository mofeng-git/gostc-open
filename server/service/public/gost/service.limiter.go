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
	var now = time.Now()
	result := cache.GetTunnelInfo(req.Client)
	if result.Code == "" || (result.ExpAt < now.Unix() && result.ChargingTye == model.GOST_CONFIG_CHARGING_CUCLE_DAY) {
		return LimiterResp{In: 1024, Out: 1024}
	}

	nodeInfo := cache.GetNodeInfo(result.NodeCode)
	if nodeInfo.LimitResetIndex != 0 && nodeInfo.LimitTotal > 0 {
		var obsUseTotal int64
		obsLimit := cache.GetNodeObsLimit(nodeInfo.Code)
		obsInfo := cache.GetTunnelObs(now.Format(time.DateOnly), result.Code)
		switch nodeInfo.LimitKind {
		case model.GOST_NODE_LIMIT_KIND_ALL:
			obsUseTotal += obsInfo.InputBytes + obsInfo.OutputBytes + obsLimit.InputBytes + obsLimit.OutputBytes
		case model.GOST_NODE_LIMIT_KIND_INPUT:
			obsUseTotal += obsInfo.InputBytes + obsLimit.InputBytes
		case model.GOST_NODE_LIMIT_KIND_OUTPUT:
			obsUseTotal += obsInfo.OutputBytes + obsLimit.OutputBytes
		}
		if obsUseTotal >= int64(nodeInfo.LimitTotal)*1024*1024*1024 {
			return LimiterResp{In: 1024, Out: 1024}
		}
	}

	speed = result.Limiter * 128 * 1024
	return LimiterResp{In: speed, Out: speed}
}
