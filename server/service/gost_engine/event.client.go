package gost_engine

import (
	"encoding/json"
	"errors"
	"github.com/lxzan/gws"
	cache2 "github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"server/global"
	"server/model"
	"server/repository"
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
	isRegister      bool

	maxRetries    int // 最大重试次数
	retryInterval time.Duration
	checkInterval time.Duration // 循环检查消息等待间隔
}

func NewClientEvent(code string, ip string, log *zap.Logger) *ClientEvent {
	return &ClientEvent{
		ip:            ip,
		code:          code,
		log:           log,
		sendAtomic:    atomic.Int32{},
		maxRetries:    10,
		retryInterval: time.Second,
		checkInterval: time.Second,
	}
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
	go event.checkRegister()
}

func (event *ClientEvent) checkRegister() {
	time.Sleep(time.Second * 10)
	if event.isRegister {
		return
	}
	if !event.isRunning {
		return
	}
	event.log.Warn("客户端连接注册超时，关闭连接", zap.String("type", "client"), zap.String("code", event.code))
	ClientStop(event.code, "客户端连接注册超时")
}

func (event *ClientEvent) register() {
	event.isRegister = true
	db, _, _ := repository.Get("")
	var hostCodes []string
	_ = db.GostClientHost.Where(db.GostClientHost.ClientCode.Eq(event.code)).Pluck(db.GostClientHost.Code, &hostCodes)
	for _, code := range hostCodes {
		ClientHostConfig(db, code)
	}
	var forwardCodes []string
	_ = db.GostClientForward.Where(db.GostClientForward.ClientCode.Eq(event.code)).Pluck(db.GostClientForward.Code, &forwardCodes)
	for _, code := range forwardCodes {
		ClientForwardConfig(db, code)
	}
	var tunnelCodes []string
	_ = db.GostClientTunnel.Where(db.GostClientTunnel.ClientCode.Eq(event.code)).Pluck(db.GostClientTunnel.Code, &tunnelCodes)
	for _, code := range tunnelCodes {
		ClientTunnelConfig(db, code)
	}
	var proxyCodes []string
	_ = db.GostClientProxy.Where(db.GostClientProxy.ClientCode.Eq(event.code)).Pluck(db.GostClientProxy.Code, &proxyCodes)
	for _, code := range proxyCodes {
		ClientProxyConfig(db, code)
	}
	var p2pCodes []string
	_ = db.GostClientP2P.Where(db.GostClientP2P.ClientCode.Eq(event.code)).Pluck(db.GostClientP2P.Code, &p2pCodes)
	for _, code := range p2pCodes {
		ClientP2PConfig(db, code)
	}
}

func (event *ClientEvent) sendLoop(socket *gws.Conn) {
	ticker := time.NewTicker(event.checkInterval)
	defer func() {
		event.log.Info("stop sendLoop",
			zap.String("type", "client"),
			zap.String("code", event.code))
		ticker.Stop()
	}()

	for event.isRunning {
		// Skip if sending is blocked
		if event.sendAtomic.Load() == 1 {
			time.Sleep(event.retryInterval)
			continue
		}

		msgChan, err := msgRegistry.PullMessage(event.code)
		if err != nil {
			time.Sleep(event.retryInterval)
			continue
		}

		select {
		case req := <-msgChan:
			event.handleMessage(socket, req)
		case <-ticker.C:
			// Continue loop
		}
	}
}

func (event *ClientEvent) handleMessage(socket *gws.Conn, req MessageData) {
	event.sendAtomic.Store(1)
	event.sendOperationId = req.OperationId
	event.WriteAny(socket, req)

	for retry := 0; retry < event.maxRetries && event.isRunning; retry++ {
		if event.sendAtomic.Load() == 0 {
			return
		}

		time.Sleep(event.retryInterval)
	}

	// Final attempt if still needed
	if event.isRunning && event.sendAtomic.Load() == 1 {
		event.WriteAny(socket, req)
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
		event.register()
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
	case "port_check":
		var data ClientPortCheckResult
		_ = msg.GetContent(&data)
		if data.Code == "" {
			return
		}
		cache.SetClientPortUse(data.Code, data.Port, data.Use, cache2.NoExpiration)
	}
}

type ClientPortCheckResult struct {
	Code string `json:"code"` // 节点编号
	Use  bool   `json:"use"`  // 1=被占用
	Port string `json:"port"` // 端口
}

func (event *ClientEvent) WriteAny(socket *gws.Conn, data any) {
	marshal, _ := json.Marshal(data)
	_ = socket.WriteString(string(marshal))
}
