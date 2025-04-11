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
	Address   string `json:"address"`
	Tls       int    `json:"tls"`
	AutoStart int    `json:"autoStart"`
	Status    int    `json:"status"`
}

func (*service) List() (result []Item) {
	for _, key := range global.ClientFS.ListKeys() {
		value, ok := global.ClientFS.Get(key)
		if !ok {
			continue
		}
		var client model.Client
		marshal, _ := json.Marshal(value)
		_ = json.Unmarshal(marshal, &client)
		result = append(result, Item{
			Key:       client.Key,
			Name:      client.Name,
			Address:   client.Address,
			Tls:       client.Tls,
			AutoStart: client.AutoStart,
			Status:    utils.TrinaryOperation(common.State.Get(key), 1, 2),
		})
	}
	return result
}
