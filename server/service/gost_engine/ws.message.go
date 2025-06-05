package gost_engine

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
	if content == nil {
		msg.Content = ""
		return
	}
	marshal, err := json.Marshal(content)
	if err != nil {
		return
	}
	msg.Content = string(marshal)
}

func NewMessage(operationId, operationType string, content any) MessageData {
	msg := MessageData{
		OperationId:   operationId,
		OperationType: operationType,
	}
	msg.SetContent(content)
	return msg
}
