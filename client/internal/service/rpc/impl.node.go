package rpcService

import (
	"errors"
	"github.com/lesismal/arpc"
	"gostc-sub/internal/common"
	"gostc-sub/internal/event_node"
	"gostc-sub/pkg/rpc_protocol/websocket"
	"net"
	"net/http"
	"time"
)

type Node struct {
	key          string
	proxyBaseUrl string
	wsUrl        string
	stopFunc     func()
}

func NewNode(wsUrl, key, proxyBaseUrl string) *Node {
	return &Node{
		key:          key,
		proxyBaseUrl: proxyBaseUrl,
		wsUrl:        wsUrl + "/rpc/node/ws",
	}
}

func (svc *Node) Start() (err error) {
	if common.State.Get(svc.key) {
		return errors.New("客户端已在运行中")
	}
	go func() {
		err = svc.run()
		for common.State.Get(svc.key) {
			if err := svc.run(); err != nil {
				common.Logger.AddLog("client", svc.key+"节点连接异常断开，等待5秒重试，err:"+err.Error())
				time.Sleep(time.Second * 5)
			}
		}
	}()
	time.Sleep(time.Second * 2)
	if err != nil {
		common.State.Set(svc.key, false)
	} else {
		common.State.Set(svc.key, true)
	}
	return err
}

func (svc *Node) Stop() {
	common.State.Set(svc.key, false)
	if svc.stopFunc != nil {
		svc.stopFunc()
	}
}

func (svc *Node) IsRunning() bool {
	return common.State.Get(svc.key)
}

func (svc *Node) run() (err error) {
	if svc.key == "" {
		return errors.New("please entry key")
	}
	client, err := arpc.NewClient(func() (net.Conn, error) {
		return websocket.Dial(svc.wsUrl, http.Header{
			"key": []string{svc.key},
		})
	})
	if err != nil {
		return err
	}
	defer client.Stop()
	client.Keepalive(time.Second * 15)
	client.Handler.SetReadTimeout(time.Second * 50)
	event, err := event_node.NewEvent(client, svc.key, svc.proxyBaseUrl)
	if err != nil {
		return err
	}
	return event.Run()
}
