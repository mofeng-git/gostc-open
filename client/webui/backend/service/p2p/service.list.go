package service

import (
	"encoding/json"
	service2 "gostc-sub/internal/service"
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
	for _, key := range global.P2PFS.ListKeys() {
		value, ok := global.P2PFS.Get(key)
		if !ok {
			continue
		}
		var p2p model.P2P
		marshal, _ := json.Marshal(value)
		_ = json.Unmarshal(marshal, &p2p)
		result = append(result, Item{
			Key:       p2p.Key,
			Name:      p2p.Name,
			Bind:      p2p.Bind,
			Port:      p2p.Port,
			Address:   p2p.Address,
			Tls:       p2p.Tls,
			AutoStart: p2p.AutoStart,
			Status:    utils.TrinaryOperation(service2.State.Get(key), 1, 2),
		})
	}
	return result
}
