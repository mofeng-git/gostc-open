package gost_engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lxzan/gws"
	"go.uber.org/zap"
	"server/global"
	"server/model"
	"server/service/common/cache"
	"sync/atomic"
	"time"
)

type NodeEvent struct {
	code            string
	log             *zap.Logger
	isRunning       bool
	sendAtomic      atomic.Int32
	sendOperationId string
}

func NewNodeEvent(code string, log *zap.Logger) *NodeEvent {
	return &NodeEvent{code: code, log: log}
}

func (event *NodeEvent) OnOpen(socket *gws.Conn) {
	event.isRunning = true
	if event.code == "" {
		event.isRunning = false
		return
	}
	cache.SetNodeOnline(event.code, true, time.Second*30)
	go func() {
		for {
			time.Sleep(time.Second * 10)
			if err := socket.WritePing(nil); err != nil {
				return
			}
		}
	}()
	go event.sendLoop(socket)
	NodeConfig(global.DB.GetDB(), event.code)
}

func (event *NodeEvent) sendLoop(socket *gws.Conn) {
	ticker := time.NewTicker(time.Second)
	defer func() {
		fmt.Println("stop sendLoop")
		ticker.Stop()
	}()
	for {
		if !event.isRunning {
			return
		}
		if event.sendAtomic.Load() == 1 {
			time.Sleep(time.Second)
			continue
		}
		msgChan, err := msgRegistry.PullMessage(event.code)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		ticker.Reset(time.Second)
		select {
		case req := <-msgChan:
			event.sendAtomic.Store(1)
			event.sendOperationId = req.OperationId
			event.WriteAny(socket, req)
			sleep := 0
			for {
				if event.sendAtomic.Load() == 0 {
					break
				} else {
					sleep++
					if sleep <= 10 {
						time.Sleep(time.Second)
						continue
					}
				}
				sleep = 0
				event.WriteAny(socket, req)
			}
		case <-ticker.C:
		}
	}
}

func (event *NodeEvent) OnClose(socket *gws.Conn, err error) {
	if !event.isRunning {
		return
	}
	event.isRunning = false
	event.log.Info("connect close")
	if err != nil && !errors.As(err, &gws.ErrConnClosed) {
		event.log.Info("connect close fail", zap.Error(err))
	}
	cache.SetNodeOnline(event.code, false, time.Second*30)
	msgRegistry.CleanMessage(event.code)
	event.sendAtomic.Store(0)
}

func (event *NodeEvent) OnPing(socket *gws.Conn, payload []byte) {
	cache.SetNodeOnline(event.code, true, time.Second*30)
	_ = socket.SetDeadline(time.Now().Add(time.Second * 30))
}

func (event *NodeEvent) OnPong(socket *gws.Conn, payload []byte) {
}

func (event *NodeEvent) OnMessage(socket *gws.Conn, message *gws.Message) {
	var msg MessageData
	_ = json.Unmarshal(message.Bytes(), &msg)
	if msg.OperationId == event.sendOperationId {
		event.sendAtomic.Store(0)
	}
	switch msg.OperationType {
	case "register":
		var data = make(map[string]any)
		_ = msg.GetContent(&data)
		cache.SetNodeVersion(event.code, data["version"].(string))
	case "logger":
		var cfg model.SystemConfigGost
		cache.GetSystemConfigGost(&cfg)
		if cfg.Logger != "1" {
			break
		}
		var logMsg = ""
		_ = msg.GetContent(&logMsg)
		var log = struct {
			Level string `json:"level"`
		}{}
		_ = json.Unmarshal([]byte(logMsg), &log)
		global.DB.GetDB().Create(&model.GostNodeLogger{
			Level:     log.Level,
			NodeCode:  event.code,
			Content:   logMsg,
			CreatedAt: time.Now().Unix(),
		})
	}
}

func (event *NodeEvent) WriteAny(socket *gws.Conn, data any) {
	marshal, _ := json.Marshal(data)
	_ = socket.WriteString(string(marshal))
}
