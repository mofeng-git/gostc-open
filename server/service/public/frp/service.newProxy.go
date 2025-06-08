package service

import (
	"errors"
	"fmt"
	"server/model"
	"server/service/common/cache"
	"strings"
	"time"
)

// {"version":"0.1.0","op":"NewProxy","content":{"user":{"user":"","metas":{"password":"2l5y63jh3r","user":"5gzxwh2ntg"},"run_id":"ade0a11d3c56d017"},"proxy_name":"bfe50ee4-72b7-495e-84a3-fe43b0a49321_proxy","proxy_type":"tcp","use_encryption":true,"use_compression":true,"bandwidth_limit":"256KB","metas":{"password":"2l5y63jh3r","user":"5gzxwh2ntg"},"remote_port":10064}}
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
	} `json:"content"`
}

func (s *service) NewProxy(req NewProxyReq) (any, error) {
	now := time.Now()
	reqCode := strings.Split(req.Content.ProxyName, "_")[0]
	tunnelCode := cache.GetGostAuth(req.Content.Metas.User, req.Content.Metas.Password)
	if reqCode != tunnelCode || tunnelCode == "" {
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

	req.Content.BandwidthLimit = fmt.Sprintf("%dKB", result.Limiter*128)
	req.Content.UseCompression = true
	req.Content.UseEncryption = true
	return req.Content, nil
}
