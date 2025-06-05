package event_client

import (
	"fmt"
	"github.com/lesismal/arpc"
	"gostc-sub/internal/common"
	"gostc-sub/pkg/utils"
	"os"
	"strconv"
	"sync"
	"time"
)

type Event struct {
	key    string
	client *arpc.Client
	svcMap *sync.Map
	tunMap *sync.Map
}

func NewEvent(client *arpc.Client, key string) (*Event, error) {
	e := Event{
		key:    key,
		client: client,
		svcMap: &sync.Map{},
		tunMap: &sync.Map{},
	}

	e.client.Handler.Handle("stop", func(c *arpc.Context) {
		var msg string
		_ = c.Bind(&msg)
		fmt.Println("停止运行原因：", msg)
		os.Exit(1)
	})
	e.client.Handler.Handle("forward_config", func(c *arpc.Context) {
		var data ClientForwardConfigData
		if err := c.Bind(&data); err != nil {
			_ = c.Write(err.Error())
			return
		}
		e.WsForwardConfig(data)
		_ = c.Write("success")
	})
	e.client.Handler.Handle("host_config", func(c *arpc.Context) {
		var data ClientHostConfigData
		if err := c.Bind(&data); err != nil {
			_ = c.Write(err.Error())
			return
		}
		e.WsHostConfig(data)
		_ = c.Write("success")
	})
	e.client.Handler.Handle("p2p_config", func(c *arpc.Context) {
		var data ClientP2PConfigData
		if err := c.Bind(&data); err != nil {
			_ = c.Write(err.Error())
			return
		}
		e.WsP2PConfig(data)
		_ = c.Write("success")
	})
	e.client.Handler.Handle("proxy_config", func(c *arpc.Context) {
		var data ClientProxyConfigData
		if err := c.Bind(&data); err != nil {
			_ = c.Write(err.Error())
			return
		}
		e.WsProxyConfig(data)
		_ = c.Write("success")
	})
	e.client.Handler.Handle("tunnel_config", func(c *arpc.Context) {
		var data ClientTunnelConfigData
		if err := c.Bind(&data); err != nil {
			_ = c.Write(err.Error())
			return
		}
		e.WsTunnelConfig(data)
		_ = c.Write("success")
	})
	e.client.Handler.Handle("remove_config", func(c *arpc.Context) {
		var data []string
		if err := c.Bind(&data); err != nil {
			_ = c.Write(err.Error())
			return
		}
		e.WsRemoveConfig(data)
		_ = c.Write("success")
	})
	e.client.Handler.Handle("port_check", func(c *arpc.Context) {
		var data = make(map[string]string)
		if err := c.Bind(&data); err != nil {
			_ = c.Write(err.Error())
			return
		}
		port, err := strconv.Atoi(data["port"])
		if err != nil {
			_ = c.Write(err.Error())
			return
		}
		_ = c.Write(utils.TrinaryOperation(utils.IsUse("0.0.0.0", port) == nil, "success", "is use"))
	})

	return &e, nil
}

func (e *Event) Run() error {
	_ = e.client.Notify("rpc/client/reg", map[string]string{
		"key":     e.key,
		"version": common.VERSION,
	}, time.Second*5)
	for {
		time.Sleep(time.Second)
		if err := e.client.CheckState(); err != nil {
			return err
		}
	}
}
