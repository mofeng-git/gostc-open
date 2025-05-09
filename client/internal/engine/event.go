package engine

import (
	"encoding/json"
	"fmt"
	"github.com/go-gost/x/registry"
	"github.com/google/uuid"
	"github.com/lxzan/gws"
	"gostc-sub/internal/common"
	"sync"
	"time"
)

type Event struct {
	server       bool
	key          string
	proxyBaseUrl string
	tp           string
	cacheBaseUrl string
	svcMap       *sync.Map
}

func NewEvent(key string, proxyBaseUrl, cacheBaseUrl string, server bool) *Event {
	tp := "client"
	if server {
		tp = "server"
	}
	return &Event{
		server:       server,
		key:          key,
		proxyBaseUrl: proxyBaseUrl,
		cacheBaseUrl: cacheBaseUrl,
		tp:           tp,
		svcMap:       &sync.Map{},
	}
}

func (e *Event) OnOpen(socket *gws.Conn) {
	common.Logger.AddLog(e.tp, "WS连接成功")
	go func() {
		for {
			time.Sleep(time.Second * 10)
			if err := socket.WritePing(nil); err != nil {
				_ = socket.WriteClose(1000, nil)
				return
			}
		}
	}()
	var regData = map[string]any{
		"version": common.VERSION,
	}
	if e.proxyBaseUrl != "" {
		regData["custom_domain"] = 1
	}
	if e.cacheBaseUrl != "" {
		regData["cache"] = 1
	}
	e.WriteAny(socket, common.NewMessage(uuid.NewString(), "register", regData))
}

func (e *Event) OnClose(socket *gws.Conn, err error) {
	common.Logger.AddLog(e.tp, "WS连接断开")
}

func (e *Event) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.SetDeadline(time.Now().Add(time.Second * 30))
}

func (e *Event) OnPong(socket *gws.Conn, payload []byte) {
}

func (e *Event) OnMessage(socket *gws.Conn, message *gws.Message) {
	var msg common.MessageData
	_ = json.Unmarshal(message.Bytes(), &msg)
	switch msg.OperationType {
	case "stop":
		var data = make(map[string]any)
		_ = msg.GetContent(&data)
		fmt.Println(data["reason"])
		common.Logger.AddLog(e.tp, data["reason"].(string))
		time.Sleep(time.Second * 2)
		common.State.Set(e.key, false)
		e.WriteAny(socket, common.NewMessage(msg.OperationId, msg.OperationType, nil))
		_ = socket.WriteClose(1000, nil)
		e.svcMap.Range(func(key, value any) bool {
			if svc := registry.ServiceRegistry().Get(key.(string)); svc != nil {
				_ = svc.Close()
			}
			return true
		})
	case "forward_config":
		var data ClientForwardConfigData
		_ = msg.GetContent(&data)
		e.WsForwardConfig(data)
		e.WriteAny(socket, common.NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "host_config":
		var data ClientHostConfigData
		_ = msg.GetContent(&data)
		e.WsHostConfig(data)
		e.WriteAny(socket, common.NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "tunnel_config":
		var data ClientTunnelConfigData
		_ = msg.GetContent(&data)
		e.WsTunnelConfig(data)
		e.WriteAny(socket, common.NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "proxy_config":
		var data ClientProxyConfigData
		_ = msg.GetContent(&data)
		e.WsProxyConfig(data)
		e.WriteAny(socket, common.NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "p2p_config":
		var data ClientP2PConfigData
		_ = msg.GetContent(&data)
		e.WsP2PConfig(data)
		e.WriteAny(socket, common.NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "remove_config":
		var names []string
		_ = msg.GetContent(&names)
		e.WsRemoveConfig(names)
		e.WriteAny(socket, common.NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "config":
		var data ConfigData
		_ = msg.GetContent(&data)
		e.WsConfig(data)
		e.WriteAny(socket, common.NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "https_domain":
		var data DomainData
		_ = msg.GetContent(&data)
		e.WsDomain(e.proxyBaseUrl, data)
		e.WriteAny(socket, common.NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "port_check":
		var data = make(map[string]string)
		_ = msg.GetContent(&data)
		e.WriteAny(socket, common.NewMessage(msg.OperationId, msg.OperationType, e.WsPortCheck(data)))
	default:
		e.WriteAny(socket, common.NewMessage(msg.OperationId, msg.OperationType, nil))
	}
}

func (e *Event) WriteAny(socket *gws.Conn, data any) {
	marshal, _ := json.Marshal(data)
	_ = socket.WriteString(string(marshal))
}
