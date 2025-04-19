package service

import (
	"encoding/json"
	"gostc-sub/internal/common"
	"gostc-sub/pkg/utils"
	"gostc-sub/webui/backend/global"
	"gostc-sub/webui/backend/model"
)

type Item struct {
	Key       string `json:"key"`
	Name      string `json:"name"`
	Bind      string `json:"bind"`
	Port      string `json:"port"`
	Address   string `json:"address"`
	Tls       int    `json:"tls"`
	AutoStart int    `json:"autoStart"`
	Status    int    `json:"status"`
}

func (*service) List() (result []Item) {
	for _, key := range global.TunnelFS.ListKeys() {
		value, ok := global.TunnelFS.Get(key)
		if !ok {
			continue
		}
		var tunnel model.Tunnel
		marshal, _ := json.Marshal(value)
		_ = json.Unmarshal(marshal, &tunnel)
		result = append(result, Item{
			Key:       tunnel.Key,
			Name:      tunnel.Name,
			Bind:      tunnel.Bind,
			Port:      tunnel.Port,
			Address:   tunnel.Address,
			Tls:       tunnel.Tls,
			AutoStart: tunnel.AutoStart,
			Status:    utils.TrinaryOperation(common.State.Get(key), 1, 2),
		})
	}
	return result
}
