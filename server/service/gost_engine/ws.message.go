package gost_engine

import (
	"encoding/json"
	"errors"
	"sync"
)

var msgRegistry = make(Message)

type Message map[string]chan MessageData

var messageLock = &sync.RWMutex{}

func (msg Message) New(code string) {
	messageLock.Lock()
	defer messageLock.Unlock()
	value := msg[code]
	if value != nil {
		close(value)
		delete(msg, code)
	}
	msg[code] = make(chan MessageData, 1000)
}

func (msg Message) PushMessage(code string, req MessageData) {
	messageLock.Lock()
	defer messageLock.Unlock()
	value := msg[code]
	if value == nil {
		return
	}
	value <- req
	msg[code] = value
}

func (msg Message) PullMessage(code string) (<-chan MessageData, error) {
	messageLock.RLock()
	defer messageLock.RUnlock()
	req := msg[code]
	if req == nil {
		return nil, errors.New("msg is nil")
	}
	return req, nil
}

func (msg Message) CleanMessage(code string) {
	messageLock.Lock()
	defer messageLock.Unlock()
	req := msg[code]
	if req == nil {
		return
	}
	close(req)
	delete(msg, code)
}

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
