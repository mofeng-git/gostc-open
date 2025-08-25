package service

import (
	"errors"
	"github.com/SianHH/frp-package/package"
	"github.com/lesismal/arpc"
	"gostc-sub/internal/common"
	"gostc-sub/internal/service/event"
	"gostc-sub/pkg/rpc_protocol/websocket"
	"net"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	key          string
	generate     common.GenerateUrl
	core         service.Service
	svcMap       *sync.Map
	svcUpdateMap *sync.Map
	stopFunc     func()
}

func NewClient(url common.GenerateUrl, key string) *Client {
	return &Client{
		key:          key,
		generate:     url,
		svcMap:       &sync.Map{},
		svcUpdateMap: &sync.Map{},
	}
}

func (svc *Client) Start() (err error) {
	if State.Get(svc.key) {
		return errors.New("客户端已在运行中")
	}
	go func() {
		State.Set(svc.key, true)
		err = svc.run()
		for State.Get(svc.key) {
			if err = svc.run(); err != nil {
				common.Logger.AddLog("client", svc.key+"客户端连接异常断开，等待5秒重试，err:"+err.Error())
				time.Sleep(time.Second * 5)
			}
		}
	}()
	time.Sleep(time.Second * 2)
	return err
}

func (svc *Client) Stop() {
	State.Set(svc.key, false)
	if svc.stopFunc != nil {
		svc.stopFunc()
	}
	svc.svcMap.Range(func(key, value any) bool {
		service.Del(key.(string))
		return true
	})
}

func (svc *Client) IsRunning() bool {
	return State.Get(svc.key)
}

func (svc *Client) run() (err error) {
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
		var reply string
		if err := client.Call("rpc/client/reg", data, &reply, time.Second*5); err != nil {
			client.Stop()
			return
		}
		if reply != "success" {
			err = errors.New(reply)
			common.Logger.AddLog("Client", reply)
			client.Stop()
			return
		}
	})

	// 注册事件
	var callback = func(key, updateTag string) {
		svc.svcMap.Store(key, true)
		svc.svcUpdateMap.Store(key, updateTag)
	}
	// true:allowUpdate false:denyUpdate
	var checkUpdate = func(key, updateTag string) bool {
		if updateTag == "" {
			return true
		}
		value, ok := svc.svcUpdateMap.Load(key)
		if !ok {
			return true
		}
		u, ok := value.(string)
		if !ok {
			return true
		}
		return u != updateTag
	}
	event.PortCheckHandle(client)
	event.HostHandle(client, callback, checkUpdate)
	event.ForwardHandle(client, callback, checkUpdate)
	event.TunnelHandle(client, callback, checkUpdate)
	event.ProxyHandle(client, callback, checkUpdate)
	event.P2PHandle(client, callback, checkUpdate)
	event.CustomCfgHandle(client, callback, checkUpdate)
	event.StopHandle(client, func() {
		svc.Stop()
	})
	event.RemoveHandle(client, func(key string) {
		svc.svcMap.Delete(key)
		svc.svcUpdateMap.Delete(key)
	})
	go svc.ping(client)
	// 开启证书目录变动监听
	certWatcherDone := certWatcher(client)
	svc.stopFunc = func() {
		certWatcherDone()
		client.Stop()
	}
	<-stopChan
	return err
}

func (svc *Client) ping(client *arpc.Client) {
	for {
		time.Sleep(time.Second * 15)
		if err := client.CheckState(); err != nil {
			if !errors.Is(err, arpc.ErrClientReconnecting) {
				return
			}
		} else {
			_ = client.CallAsync("rpc/client/ping", nil, func(c *arpc.Context, err error) {}, time.Second*5)
		}
	}
}
