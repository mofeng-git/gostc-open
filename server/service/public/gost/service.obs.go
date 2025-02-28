package service

import (
	"server/pkg/utils"
	"server/service/common/cache"
	"time"
)

type ObsReq struct {
	Events []struct {
		Kind    string `json:"kind"`
		Service string `json:"service"`
		Client  string `json:"client"`
		Type    string `json:"type"`
		Stats   struct {
			TotalConns   int   `json:"totalConns"`
			CurrentConns int   `json:"currentConns"`
			InputBytes   int64 `json:"inputBytes"`
			OutputBytes  int64 `json:"outputBytes"`
			TotalErrs    int   `json:"totalErrs"`
		} `json:"stats"`
	} `json:"events"`
}

type ObsResp struct {
	Ok bool `json:"ok"`
}

func (service *service) Obs(tunnel string, req ObsReq) ObsResp {
	dateOnly := time.Now().Format(time.DateOnly)
	for _, event := range req.Events {
		inputBytes := event.Stats.InputBytes
		outputBytes := event.Stats.OutputBytes
		if tunnel == "" {
			// 反转流量记录数据
			// inputBytes 对内网服务发送的数据流量
			// outputBytes 内网服务对外发送的数据流量
			inputBytes = event.Stats.OutputBytes
			outputBytes = event.Stats.InputBytes
		}
		if event.Type == "status" || (inputBytes == 0 && outputBytes == 0) {
			continue
		}
		tunnelCode := utils.TrinaryOperation(tunnel == "", event.Client, tunnel)
		tunnelInfo := cache.GetTunnelInfo(tunnelCode)
		nodeVersion := cache.GetNodeVersion(tunnelInfo.NodeCode)
		if nodeVersion > "v1.1.2" && tunnel != "" {
			continue
		}
		go cache.IncreaseObs(dateOnly, tunnelInfo.Code, tunnelInfo.ClientCode, tunnelInfo.NodeCode, tunnelInfo.UserCode, cache.TunnelObs{
			InputBytes:  inputBytes,
			OutputBytes: outputBytes,
		})
	}
	return ObsResp{Ok: true}
}
