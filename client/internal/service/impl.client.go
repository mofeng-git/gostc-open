package service

import (
	"crypto/tls"
	"errors"
	"github.com/lxzan/gws"
	"gostc-sub/internal/common"
	"gostc-sub/internal/engine"
	"net/http"
	"time"
)

type Client struct {
	key   string
	wsUrl string
	svc   *gws.Conn
}

func NewClient(wsUrl, key string) *Client {
	return &Client{
		key:   key,
		wsUrl: wsUrl + "/api/v1/public/gost/client/ws",
	}
}

func (svc *Client) Start() (err error) {
	if common.State.Get(svc.key) {
		return errors.New("客户端已在运行中")
	}
	go func() {
		err = svc.run()
		for {
			if !common.State.Get(svc.key) {
				return
			}
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
	if svc.svc != nil {
		_ = svc.svc.WriteClose(1000, nil)
	}
}

func (svc *Client) IsRunning() bool {
	return common.State.Get(svc.key)
}

func (svc *Client) run() (err error) {
	if svc.key == "" {
		return errors.New("please entry key")
	}
	svc.svc, _, err = gws.NewClient(engine.NewEvent(svc.key, "", false), &gws.ClientOption{
		Addr:      svc.wsUrl,
		TlsConfig: &tls.Config{InsecureSkipVerify: true},
		NewDialer: func() (gws.Dialer, error) {
			return common.NewDialer(), nil
		},
		RequestHeader: http.Header{"key": []string{svc.key}},
		PermessageDeflate: gws.PermessageDeflate{
			Enabled: true,
		},
	})
	if err != nil {
		common.Logger.AddLog("client", svc.key+"客户端启动失败，err:"+err.Error())
		return err
	}
	svc.svc.ReadLoop()
	return nil
}
