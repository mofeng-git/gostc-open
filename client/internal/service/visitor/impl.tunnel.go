package service

import (
	"errors"
	"github.com/SianHH/frp-package/package"
	"github.com/SianHH/frp-package/package/frpc"
	v1 "github.com/SianHH/frp-package/pkg/config/v1"
	"gopkg.in/yaml.v3"
	"gostc-sub/internal/common"
	service2 "gostc-sub/internal/service"
	"io"
	"net/http"
	"os"
	"time"
)

type Tunnel struct {
	key      string
	generate common.GenerateUrl
	bindPort int
	bindAddr string
	core     service.Service
	stopFunc func()
}

func NewTunnel(url common.GenerateUrl, key, bindAddr string, bindPort int) *Tunnel {
	return &Tunnel{
		key:      key,
		generate: url,
		bindPort: bindPort,
		bindAddr: bindAddr,
		core:     nil,
	}
}

func (svc *Tunnel) Start() (err error) {
	if service2.State.Get(svc.key) {
		return errors.New("私有隧道已在运行中")
	}
	go func() {
		err = svc.run()
		for service2.State.Get(svc.key) {
			if err := svc.run(); err != nil {
				common.Logger.AddLog("client", svc.key+"私有隧道异常停止，等待5秒重试，err:"+err.Error())
				time.Sleep(time.Second * 5)
			}
		}
	}()
	time.Sleep(time.Second * 2)
	if err != nil {
		service2.State.Set(svc.key, false)
	} else {
		service2.State.Set(svc.key, true)
	}
	return err
}

func (svc *Tunnel) Stop() {
	service2.State.Set(svc.key, false)
	if svc.stopFunc != nil {
		svc.stopFunc()
	}
	svc.stopFunc = nil
}

func (svc *Tunnel) IsRunning() bool {
	return service2.State.Get(svc.key)
}

func (svc *Tunnel) run() (err error) {
	if svc.key == "" {
		return errors.New("please entry key")
	}
	config, err := svc.loadConfig(svc.generate.HttpUrl() + "/api/v1/public/frp/visitorTunnel?key=" + svc.key)
	if err != nil {
		return err
	}
	var visitorCfgs []v1.VisitorConfigurer
	if config.STCP.Name != "" {
		config.STCP.BindAddr = svc.bindAddr
		config.STCP.BindPort = svc.bindPort
		config.STCP.Plugin = v1.TypedVisitorPluginOptions{}
		visitorCfgs = append(visitorCfgs, &config.STCP)
	}
	if config.SUDP.Name != "" {
		config.SUDP.BindAddr = svc.bindAddr
		config.SUDP.BindPort = svc.bindPort
		config.SUDP.Plugin = v1.TypedVisitorPluginOptions{}
		visitorCfgs = append(visitorCfgs, &config.SUDP)
	}
	config.Common.Transport.ProxyURL = os.Getenv("GOSTC_TRANSPORT_PROXY_URL")
	svc.core, err = frpc.NewService(config.Common, nil, visitorCfgs)
	if err != nil {
		return err
	}
	if err := svc.core.Start(); err != nil {
		return err
	}
	svc.stopFunc = func() {
		svc.core.Stop()
	}
	svc.core.Wait()
	return nil
}

func (svc *Tunnel) loadConfig(url string) (result TunnelConfig, err error) {
	response, err := http.Get(url)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return result, err
	}
	err = yaml.Unmarshal(bytes, &result)
	return result, err
}
