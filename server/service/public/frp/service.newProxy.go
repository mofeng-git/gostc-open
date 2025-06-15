package service

import (
	"errors"
	"fmt"
	"server/model"
	cache2 "server/repository/cache"
	"strings"
	"time"
)

type NewProxyReq struct {
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
		ProxyName      string `json:"proxy_name"`
		ProxyType      string `json:"proxy_type"`
		UseEncryption  bool   `json:"use_encryption"`
		UseCompression bool   `json:"use_compression"`
		BandwidthLimit string `json:"bandwidth_limit"`
		Metas          struct {
			Password string `json:"password"`
			User     string `json:"user"`
		} `json:"metas"`
		RemotePort        int      `json:"remote_port"`
		CustomDomains     []string `json:"custom_domains"`
		HostHeaderRewrite string   `json:"host_header_rewrite"`
		Headers           struct {
			XFromWhere string `json:"x-from-where"`
		} `json:"headers"`
		Sk         string   `json:"sk"`
		AllowUsers []string `json:"allow_users"`
	} `json:"content"`
}

func (s *service) NewProxy(req NewProxyReq) (any, error) {
	now := time.Now()
	reqCode := strings.Split(req.Content.ProxyName, "_")[0]
	tunnelCode := cache2.GetGostAuth(req.Content.Metas.User, req.Content.Metas.Password)
	if reqCode != tunnelCode || tunnelCode == "" {
		return nil, errors.New("unauthorized")
	}
	result := cache2.GetTunnelInfo(tunnelCode)
	if result.Code == "" || (result.ExpAt < now.Unix() && result.ChargingTye == model.GOST_CONFIG_CHARGING_CUCLE_DAY) {
		return nil, errors.New("expired")
	}

	nodeInfo := cache2.GetNodeInfo(result.NodeCode)
	if nodeInfo.LimitResetIndex != 0 && nodeInfo.LimitTotal > 0 {
		var obsUseTotal int64
		obsLimit := cache2.GetNodeObsLimit(nodeInfo.Code)
		obsInfo := cache2.GetTunnelObs(now.Format(time.DateOnly), result.Code)
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

	req.Content.BandwidthLimit = fmt.Sprintf("%dKB", result.Limiter*128)
	req.Content.UseCompression = true
	req.Content.UseEncryption = true
	return req.Content, nil
}
