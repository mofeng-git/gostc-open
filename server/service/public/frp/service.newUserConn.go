package service

import (
	"errors"
	"server/model"
	"server/service/common/cache"
	"strings"
	"time"
)

// {"version":"0.1.0","op":"NewUserConn","content":{"user":{"user":"","metas":{"password":"2l5y63jh3r","user":"5gzxwh2ntg"},"run_id":"ade0a11d3c56d017"},"proxy_name":"bfe50ee4-72b7-495e-84a3-fe43b0a49321_proxy", "proxy_type":"tcp","remote_addr":"127.0.0.1:8151"}}
type NewUserConnReq struct {
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
		ProxyName  string `json:"proxy_name"`
		ProxyType  string `json:"proxy_type"`
		RemoteAddr string `json:"remote_addr"`
	} `json:"content"`
}

func (s *service) NewUserConn(req NewUserConnReq) (any, error) {
	now := time.Now()
	reqCode := strings.Split(req.Content.ProxyName, "_")[0]
	tunnelCode := cache.GetGostAuth(req.Content.User.Metas.User, req.Content.User.Metas.Password)
	if tunnelCode == "" || reqCode != tunnelCode {
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
