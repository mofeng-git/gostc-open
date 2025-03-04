package common

import (
	"encoding/json"
)

type MessageData struct {
	OperationId   string `json:"operationId"`
	OperationType string `json:"operationType"`
	Content       string `json:"content"`
}

func (msg *MessageData) GetContent(data any) error {
	return json.Unmarshal([]byte(msg.Content), data)
}

func (msg *MessageData) SetContent(content any) {
	data := ""
	if content != nil {
		marshal, _ := json.Marshal(content)
		data = string(marshal)
	}
	msg.Content = data
}

func NewMessage(operationId, operationType string, content any) MessageData {
	data := ""
	if content != nil {
		marshal, _ := json.Marshal(content)
		data = string(marshal)
	}
	return MessageData{
		OperationId:   operationId,
		OperationType: operationType,
		Content:       data,
	}
}
