package service

import (
	"errors"
	"fmt"
	"server/model"
	"server/service/common/cache"
	"strings"
	"time"
)

type NewReq struct {
	Version string `json:"version"`
	Op      string `json:"op"`
	Content struct {
		User struct {
			User  string      `json:"user"`
			Metas interface{} `json:"metas"`
			RunId string      `json:"run_id"`
		} `json:"user"`
		ProxyName      string   `json:"proxy_name"`
		ProxyType      string   `json:"proxy_type"`
		UseCompression bool     `json:"use_compression"`
		BandwidthLimit string   `json:"bandwidth_limit"`
		Sk             string   `json:"sk"`
		AllowUsers     []string `json:"allow_users"`
	} `json:"content"`
}

func (s *service) New(req NewReq) (any, bool, error) {
	switch req.Content.ProxyType {
	case "stcp":
		tunnelCode := strings.ReplaceAll(req.Content.ProxyName, "stcp_", "")
		var speed = 1
		result := cache.GetTunnelInfo(tunnelCode)
		if result.Code == "" || (result.ExpAt < time.Now().Unix() && result.ChargingTye == model.GOST_CONFIG_CHARGING_CUCLE_DAY) {
			speed = 1
		} else {
			speed = result.Limiter * 128
		}
		req.Content.BandwidthLimit = fmt.Sprintf("%dKB", speed)
		return req.Content, true, nil
	case "xtcp":
		return nil, false, nil
	default:
		return nil, false, errors.New("未知类型")
	}
}
