package service

import (
	"gostc-sub/internal/common"
	"time"
)

type Item struct {
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Message   string `json:"message"`
}

func (*service) List() (result []Item) {
	for _, item := range common.Logger.GetLogs() {
		result = append(result, Item{
			Timestamp: item.Timestamp.Format(time.DateTime),
			Type:      item.Type,
			Message:   item.Message,
		})
	}
	return result
}
