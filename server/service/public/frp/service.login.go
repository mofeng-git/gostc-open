package service

import (
	"errors"
	"server/model"
	"server/service/common/cache"
	"time"
)

// {"version":"0.1.0","op":"Login","content":{"version":"0.63.0","os":"windows","arch":"amd64","privilege_key":"2932ca6e85efa28d5913b4260dd0049b","timestamp":1749310313,"metas":{"password":"2l5y63jh3r","user":"5gzxwh2ntg"},"client_spec":{},"pool_count":1,"client_address":"127.0.0.1:8145"}}

type LoginReq struct {
	Version string `json:"version"`
	Op      string `json:"op"`
	Content struct {
		Version      string `json:"version"`
		Os           string `json:"os"`
		Arch         string `json:"arch"`
		PrivilegeKey string `json:"privilege_key"`
		Timestamp    int    `json:"timestamp"`
		Metas        struct {
			Password string `json:"password"`
			User     string `json:"user"`
		} `json:"metas"`
		ClientSpec struct {
		} `json:"client_spec"`
		PoolCount     int    `json:"pool_count"`
		ClientAddress string `json:"client_address"`
	} `json:"content"`
}

func (s *service) Login(req LoginReq) (any, error) {
	now := time.Now()
	code := cache.GetGostAuth(req.Content.Metas.User, req.Content.Metas.Password)
	if code == "" {
		return nil, errors.New("unauthorized")
	}
	result := cache.GetTunnelInfo(code)
	if result.Code == "" || (result.ExpAt < now.Unix() && result.ChargingTye == model.GOST_CONFIG_CHARGING_CUCLE_DAY) {
		return nil, errors.New("expired")
	}
	return nil, nil
}
