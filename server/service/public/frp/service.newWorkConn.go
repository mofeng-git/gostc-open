package service

import (
	"errors"
	"server/model"
	"server/service/common/cache"
	"time"
)

// {"version":"0.1.0","op":"NewWorkConn","content":{"user":{"user":"","metas":{"password":"4i16bw02sf","user":"gqpudmfnf9"},"run_id":"1df8d4ea578e17a3"},"run_id":"1df8d4ea578e17a3"}}
type NewWorkConnReq struct {
	Version string `json:"version"`
	Op      string `json:"op"`
	Content struct {
		User struct {
			User  string `json:"user"`
			Metas struct {
				Password string `json:"password"`
				User     string `json:"user"`
			} `json:"metas"`
			RunId string `json:"run_id"`
		} `json:"user"`
		RunId string `json:"run_id"`
	} `json:"content"`
}

func (s *service) NewWorkConn(req NewWorkConnReq) (any, error) {
	now := time.Now()
	tunnelCode := cache.GetGostAuth(req.Content.User.Metas.User, req.Content.User.Metas.Password)
	if tunnelCode == "" {
		return nil, errors.New("unauthorized")
	}
	result := cache.GetTunnelInfo(tunnelCode)
	if result.Code == "" || (result.ExpAt < now.Unix() && result.ChargingTye == model.GOST_CONFIG_CHARGING_CUCLE_DAY) {
		return nil, errors.New("expired")
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
			return nil, errors.New("insufficient node traffic")
		}
	}
	return nil, nil
}
