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
	"time"
)

type P2P struct {
	key      string
	httpUrl  string
	bindPort int
	bindAddr string
	core     service.Service
}

func NewP2P(httpUrl string, key, bindAddr string, bindPort int) *P2P {
	return &P2P{
		key:      key,
		httpUrl:  httpUrl + "/api/v1/public/frp/visitorP2P?key=" + key,
		bindPort: bindPort,
		bindAddr: bindAddr,
	}
}

func (svc *P2P) Start() (err error) {
	if service2.State.Get(svc.key) {
		return errors.New("P2P隧道已在运行中")
	}
	go func() {
		err = svc.run()
		for service2.State.Get(svc.key) {
			if err := svc.run(); err != nil {
				common.Logger.AddLog("client", svc.key+"P2P隧道异常停止，等待5秒重试，err:"+err.Error())
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

func (svc *P2P) Stop() {
	service2.State.Set(svc.key, false)
	if svc.core != nil {
		svc.core.Stop()
	}
}

func (svc *P2P) IsRunning() bool {
	return service2.State.Get(svc.key)
}

func (svc *P2P) run() (err error) {
	if svc.key == "" {
		return errors.New("please entry key")
	}
	config, err := svc.loadConfig(svc.httpUrl)
	if err != nil {
		return err
	}
	var visitorCfgs []v1.VisitorConfigurer
	if config.STCP.Name != "" {
		config.STCP.Plugin = v1.TypedVisitorPluginOptions{}
		visitorCfgs = append(visitorCfgs, &config.STCP)
	}
	if config.XTCP.Name != "" {
		config.XTCP.BindAddr = svc.bindAddr
		config.XTCP.BindPort = svc.bindPort
		config.XTCP.Plugin = v1.TypedVisitorPluginOptions{}
		visitorCfgs = append(visitorCfgs, &config.XTCP)
	}
	svc.core, err = frpc.NewService(config.Common, nil, visitorCfgs)
	if err != nil {
		return err
	}
	if err := svc.core.Start(); err != nil {
		return err
	}
	svc.core.Wait()
	return nil
}

func (svc *P2P) loadConfig(url string) (result P2PConfig, err error) {
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
