package gost_engine

import (
	"encoding/json"
	"errors"
	"github.com/lxzan/gws"
	"go.uber.org/zap"
	"server/global"
	"server/model"
	"server/service/common/cache"
	"sync/atomic"
	"time"
)

type ClientEvent struct {
	ip              string
	code            string
	log             *zap.Logger
	isRunning       bool
	sendAtomic      atomic.Int32
	sendOperationId string
}

func NewClientEvent(code string, ip string, log *zap.Logger) *ClientEvent {
	return &ClientEvent{code: code, log: log, ip: ip}
}

func (event *ClientEvent) OnOpen(socket *gws.Conn) {
	event.isRunning = true
	if event.code == "" {
		event.isRunning = false
		return
	}
	go func() {
		for {
			time.Sleep(time.Second * 10)
			if err := socket.WritePing(nil); err != nil {
				return
			}
		}
	}()
	cache.SetClientOnline(event.code, true, time.Second*30)
	cache.SetClientLastTime(event.code)
	cache.SetClientIp(event.code, event.ip)

	go event.sendLoop(socket)
	var hostCodes []string
	global.DB.GetDB().Model(&model.GostClientHost{}).Where("client_code = ?", event.code).Pluck("code", &hostCodes)
	for _, code := range hostCodes {
		ClientHostConfig(global.DB.GetDB(), code)
	}
	var forwardCodes []string
	global.DB.GetDB().Model(&model.GostClientForward{}).Where("client_code = ?", event.code).Pluck("code", &forwardCodes)
	for _, code := range forwardCodes {
		ClientForwardConfig(global.DB.GetDB(), code)
	}
	var tunnelCodes []string
	global.DB.GetDB().Model(&model.GostClientTunnel{}).Where("client_code = ?", event.code).Pluck("code", &tunnelCodes)
	for _, code := range tunnelCodes {
		ClientTunnelConfig(global.DB.GetDB(), code)
	}
}

func (event *ClientEvent) sendLoop(socket *gws.Conn) {
	ticker := time.NewTicker(time.Second)
	defer func() {
		event.log.Info("stop sendLoop", zap.String("code", event.code))
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
		select {
		case req := <-msgChan:
			event.sendAtomic.Store(1)
			event.sendOperationId = req.OperationId
			event.WriteAny(socket, req)
			sleep := 0
			for {
				if !event.isRunning {
					return
				}
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

func (event *ClientEvent) OnClose(socket *gws.Conn, err error) {
	if !event.isRunning {
		return
	}
	event.isRunning = false
	event.log.Info("connect close")
	if err != nil && !errors.As(err, &gws.ErrConnClosed) {
		event.log.Info("connect close fail", zap.Error(err))
	}
	cache.SetClientOnline(event.code, false, time.Second*30)
	msgRegistry.CleanMessage(event.code)
	event.sendAtomic.Store(0)
}

func (event *ClientEvent) OnPing(socket *gws.Conn, payload []byte) {
	cache.SetClientOnline(event.code, true, time.Second*30)
	_ = socket.SetDeadline(time.Now().Add(time.Second * 30))
}

func (event *ClientEvent) OnPong(socket *gws.Conn, payload []byte) {

}

func (event *ClientEvent) OnMessage(socket *gws.Conn, message *gws.Message) {
	var msg MessageData
	_ = json.Unmarshal(message.Bytes(), &msg)
	if msg.OperationId == event.sendOperationId {
		event.sendAtomic.Store(0)
	}
	switch msg.OperationType {
	case "register":
		var data = make(map[string]any)
		_ = msg.GetContent(&data)
		cache.SetClientVersion(event.code, data["version"].(string))
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
		global.DB.GetDB().Create(&model.GostClientLogger{
			Level:      log.Level,
			ClientCode: event.code,
			Content:    logMsg,
			CreatedAt:  time.Now().Unix(),
		})
	}
	//fmt.Println(string(message.Bytes()))
}

func (event *ClientEvent) WriteAny(socket *gws.Conn, data any) {
	marshal, _ := json.Marshal(data)
	_ = socket.WriteString(string(marshal))
}
