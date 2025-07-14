package service

import (
	"errors"
	service "github.com/SianHH/frp-package/package"
	"github.com/lesismal/arpc"
	"gostc-sub/internal/common"
	"gostc-sub/internal/service/event"
	"gostc-sub/pkg/rpc_protocol/websocket"
	"net"
	"net/http"
	"sync"
	"time"
)

type Node struct {
	key          string
	proxyBaseUrl string
	generate     common.GenerateUrl
	svcMap       *sync.Map
	stopFunc     func()
}

func NewNode(url common.GenerateUrl, key, proxyBaseUrl string) *Node {
	return &Node{
		key:          key,
		proxyBaseUrl: proxyBaseUrl,
		svcMap:       &sync.Map{},
		generate:     url,
	}
}

func (svc *Node) Start() (err error) {
	if State.Get(svc.key) {
		return errors.New("客户端已在运行中")
	}
	go func() {
		State.Set(svc.key, true)
		err = svc.run()
		for State.Get(svc.key) {
			if err = svc.run(); err != nil {
				common.Logger.AddLog("client", svc.key+"节点连接异常断开，等待5秒重试，err:"+err.Error())
				time.Sleep(time.Second * 5)
			}
		}
	}()
	time.Sleep(time.Second * 2)
	return err
}

func (svc *Node) Stop() {
	State.Set(svc.key, false)
	if svc.stopFunc != nil {
		svc.stopFunc()
	}
	svc.svcMap.Range(func(key, value any) bool {
		service.Del(key.(string))
		return true
	})
}

func (svc *Node) IsRunning() bool {
	return State.Get(svc.key)
}

func (svc *Node) run() (err error) {
	if svc.key == "" {
		return errors.New("please entry key")
	}
	client, err := arpc.NewClient(func() (net.Conn, error) {
		return websocket.Dial(svc.generate.WsUrl()+"/rpc/ws", http.Header{
			"key": []string{svc.key},
		})
	})
	if err != nil {
		return err
	}
	defer client.Stop()
	//client.Keepalive(time.Second * 15)
	client.Handler.SetReadTimeout(time.Second * 50)
	var stopChan = make(chan struct{})
	client.Handler.HandleDisconnected(func(client *arpc.Client) {
		stopChan <- struct{}{}
	})

	// 连接成功，发送注册请求
	client.Handler.HandleConnected(func(client *arpc.Client) {
		var data = map[string]string{
			"key":     svc.key,
			"version": common.VERSION,
		}
		if svc.proxyBaseUrl != "" {
			data["domain"] = "1"
		}
		var reply string
		if err := client.Call("rpc/node/reg", data, &reply, time.Second*5); err != nil {
			client.Stop()
			return
		}
		if reply != "success" {
			err = errors.New(reply)
			common.Logger.AddLog("Node", reply)
			client.Stop()
		}
	})

	// 注册事件
	if err := InitMetrics(client); err != nil {
		return err
	}

	var callback = func(key string) {
		svc.svcMap.Store(key, true)
	}
	event.ServerHandle(client, svc.generate.HttpUrl(), callback)
	event.PortCheckHandle(client)
	event.ServerDomainHandle(client, svc.proxyBaseUrl)
	go svc.ping(client)
	event.StopHandle(client, func() {
		svc.Stop()
	})
	<-stopChan
	return err
}

func (svc *Node) ping(client *arpc.Client) {
	for {
		time.Sleep(time.Second * 15)
		if err := client.CheckState(); err != nil {
			if !errors.Is(err, arpc.ErrClientReconnecting) {
				return
			}
		} else {
			_ = client.CallAsync("rpc/node/ping", nil, func(c *arpc.Context, err error) {}, time.Second*5)
		}
	}
}
