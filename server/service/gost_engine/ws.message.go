package gost_engine

import (
	"encoding/json"
	"fmt"
	"sync"
)

var msgRegistry = make(Message)

type Message map[string]chan MessageData

var messageLock = &sync.RWMutex{}

func (msg Message) New(code string) {
	messageLock.Lock()
	defer messageLock.Unlock()
	if ch, exists := msg[code]; exists {
		close(ch)
	}
	msg[code] = make(chan MessageData, 1000)
}

func (msg Message) PushMessage(code string, req MessageData) error {
	messageLock.RLock()
	ch, exists := msg[code]
	messageLock.RUnlock()

	if !exists {
		return fmt.Errorf("channel for code %s does not exist", code)
	}

	select {
	case ch <- req:
		return nil
	default:
		return fmt.Errorf("channel for code %s is full", code)
	}
}

func (msg Message) PullMessage(code string) (<-chan MessageData, error) {
	messageLock.RLock()
	defer messageLock.RUnlock()

	if ch, exists := msg[code]; exists {
		return ch, nil
	}
	return nil, fmt.Errorf("channel for code %s does not exist", code)
}

func (msg Message) CleanMessage(code string) {
	messageLock.Lock()
	defer messageLock.Unlock()

	if ch, exists := msg[code]; exists {
		close(ch)
		delete(msg, code)
	}
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
