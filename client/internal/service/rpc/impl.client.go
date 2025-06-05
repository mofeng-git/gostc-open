package rpcService

import (
	"errors"
	"github.com/lesismal/arpc"
	"gostc-sub/internal/common"
	"gostc-sub/internal/event_client"
	"gostc-sub/pkg/rpc_protocol/websocket"
	"net"
	"net/http"
	"time"
)

type Client struct {
	key      string
	wsUrl    string
	stopFunc func()
}

func NewClient(wsUrl, key string) *Client {
	return &Client{
		key:   key,
		wsUrl: wsUrl + "/rpc/client/ws",
	}
}

func (svc *Client) Start() (err error) {
	if common.State.Get(svc.key) {
		return errors.New("客户端已在运行中")
	}
	go func() {
		err = svc.run()
		for common.State.Get(svc.key) {
			if err := svc.run(); err != nil {
				common.Logger.AddLog("client", svc.key+"客户端连接异常断开，等待5秒重试，err:"+err.Error())
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

func (svc *Client) Stop() {
	common.State.Set(svc.key, false)
	if svc.stopFunc != nil {
		svc.stopFunc()
	}
}

func (svc *Client) IsRunning() bool {
	return common.State.Get(svc.key)
}

func (svc *Client) run() (err error) {
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
	event, err := event_client.NewEvent(client, svc.key)
	if err != nil {
		return err
	}
	return event.Run()
}
