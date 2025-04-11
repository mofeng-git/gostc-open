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

type Server struct {
	key          string
	proxyBaseUrl string
	wsUrl        string
	svc          *gws.Conn
	event        *engine.Event
}

func NewServer(wsUrl, key, proxyBaseUrl string) *Server {
	return &Server{
		key:          key,
		proxyBaseUrl: proxyBaseUrl,
		wsUrl:        wsUrl + "/api/v1/public/gost/node/ws",
	}
}

func (svc *Server) Start() (err error) {
	if common.State.Get(svc.key) {
		return errors.New("服务端已在运行中")
	}
	go func() {
		err = svc.run()
		for {
			if !common.State.Get(svc.key) {
				return
			}
			if err := svc.run(); err != nil {
				common.Logger.AddLog("server", svc.key+"节点连接异常断开，等待5秒重试，err:"+err.Error())
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

func (svc *Server) Stop() {
	common.State.Set(svc.key, false)
	if svc.svc != nil {
		_ = svc.svc.WriteClose(1000, nil)
	}
}

func (svc *Server) IsRunning() bool {
	return common.State.Get(svc.key)
}

func (svc *Server) run() (err error) {
	if svc.key == "" {
		return errors.New("please entry key")
	}
	svc.event = engine.NewEvent(svc.key, svc.proxyBaseUrl, true)
	svc.svc, _, err = gws.NewClient(svc.event, &gws.ClientOption{
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
		common.Logger.AddLog("server", svc.key+"节点启动失败，err:"+err.Error())
		return err
	}
	svc.svc.ReadLoop()
	return nil
}
