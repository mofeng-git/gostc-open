package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lxzan/gws"
	"os"
	"time"
)

var SvcMap = make(map[string]bool)

type Event struct {
	server       bool
	key          string
	proxyBaseUrl string
}

func NewEvent(key string, proxyBaseUrl string, server bool) *Event {
	return &Event{
		server:       server,
		key:          key,
		proxyBaseUrl: proxyBaseUrl,
	}
}

func (e *Event) OnOpen(socket *gws.Conn) {
	fmt.Println("connect success")
	go func() {
		for {
			time.Sleep(time.Second * 10)
			if err := socket.WritePing(nil); err != nil {
				fmt.Println("send ping msg fail", err)
				_ = socket.WriteClose(1000, nil)
				return
			}
		}
	}()
	version := CLIENT_VERSION
	if e.server {
		version = SERVER_VERSION
	}
	var regData = map[string]any{
		"version": version,
	}
	if e.proxyBaseUrl != "" {
		regData["custom_domain"] = 1
	}
	e.WriteAny(socket, NewMessage(uuid.NewString(), "register", regData))
}

func (e *Event) OnClose(socket *gws.Conn, err error) {
	if !errors.Is(err, gws.ErrConnClosed) {
		fmt.Println("conn close", err)
	}
}

func (e *Event) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.SetDeadline(time.Now().Add(time.Second * 30))
}

func (e *Event) OnPong(socket *gws.Conn, payload []byte) {
}

func (e *Event) OnMessage(socket *gws.Conn, message *gws.Message) {
	var msg MessageData
	_ = json.Unmarshal(message.Bytes(), &msg)
	switch msg.OperationType {
	case "stop":
		var data = make(map[string]any)
		_ = msg.GetContent(&data)
		fmt.Println(data["reason"])
		e.WriteAny(socket, NewMessage(msg.OperationId, msg.OperationType, nil))
		time.Sleep(time.Second * 2)
		os.Exit(0)
	case "forward_config":
		var data ClientForwardConfigData
		_ = msg.GetContent(&data)
		WsForwardConfig(data)
		e.WriteAny(socket, NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "host_config":
		var data ClientHostConfigData
		_ = msg.GetContent(&data)
		WsHostConfig(data)
		e.WriteAny(socket, NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "tunnel_config":
		var data ClientTunnelConfigData
		_ = msg.GetContent(&data)
		WsTunnelConfig(data)
		e.WriteAny(socket, NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "proxy_config":
		var data ClientProxyConfigData
		_ = msg.GetContent(&data)
		WsProxyConfig(data)
		e.WriteAny(socket, NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "p2p_config":
		var data ClientP2PConfigData
		_ = msg.GetContent(&data)
		WsP2PConfig(data)
		e.WriteAny(socket, NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "remove_config":
		var names []string
		_ = msg.GetContent(&names)
		WsRemoveConfig(names)
		e.WriteAny(socket, NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "config":
		var data ConfigData
		_ = msg.GetContent(&data)
		WsConfig(data)
		e.WriteAny(socket, NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "https_domain":
		var data HttpsDomainData
		_ = msg.GetContent(&data)
		WsHttpsDomain(e.proxyBaseUrl, data)
		e.WriteAny(socket, NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "port_check":
		var data = make(map[string]string)
		_ = msg.GetContent(&data)
		e.WriteAny(socket, NewMessage(msg.OperationId, msg.OperationType, WsPortCheck(data)))
	default:
		e.WriteAny(socket, NewMessage(msg.OperationId, msg.OperationType, nil))
	}
}

func (e *Event) WriteAny(socket *gws.Conn, data any) {
	marshal, _ := json.Marshal(data)
	_ = socket.WriteString(string(marshal))
}
