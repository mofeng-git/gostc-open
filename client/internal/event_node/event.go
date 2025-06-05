package event_node

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
	key          string
	client       *arpc.Client
	proxyBaseUrl string
	svcMap       *sync.Map
	tunMap       *sync.Map
}

func NewEvent(client *arpc.Client, key, proxyBaseUrl string) (*Event, error) {
	e := Event{
		key:          key,
		client:       client,
		proxyBaseUrl: proxyBaseUrl,
		svcMap:       &sync.Map{},
		tunMap:       &sync.Map{},
	}

	e.client.Handler.Handle("stop", func(c *arpc.Context) {
		var msg string
		_ = c.Bind(&msg)
		fmt.Println("停止运行原因：", msg)
		os.Exit(1)
	})

	e.client.Handler.Handle("config", func(c *arpc.Context) {
		var data ConfigData
		if err := c.Bind(&data); err != nil {
			_ = c.Write(err.Error())
			return
		}
		e.WsConfig(data)
		_ = c.Write("success")
	})
	e.client.Handler.Handle("https_domain", func(c *arpc.Context) {
		var data DomainData
		if err := c.Bind(&data); err != nil {
			_ = c.Write(err.Error())
			return
		}
		e.WsDomain(e.proxyBaseUrl, data)
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
	var data = map[string]string{
		"key":     e.key,
		"version": common.VERSION,
	}
	if e.proxyBaseUrl != "" {
		data["custom_domain"] = "1"
	}
	_ = e.client.Notify("rpc/node/reg", data, time.Second*5)
	for {
		time.Sleep(time.Second)
		if err := e.client.CheckState(); err != nil {
			return err
		}
	}
}
