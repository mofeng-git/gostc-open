package service

import (
	"errors"
	"gostc-sub/internal/common"
	"gostc-sub/pkg/p2p/frpc"
	v1 "gostc-sub/pkg/p2p/pkg/config/v1"
	"gostc-sub/pkg/p2p/registry"
	"gostc-sub/pkg/utils"
	"strconv"
)

type P2P struct {
	key     string
	httpUrl string
	bind    string
	port    string
}

func NewP2P(httpUrl, key, bind, port string) *P2P {
	return &P2P{
		key:     key,
		httpUrl: httpUrl,
		bind:    bind,
		port:    port,
	}
}

func (svc *P2P) Start() (err error) {
	if !utils.ValidatePort(svc.port) {
		return errors.New("本地端口格式错误")
	}
	port, err := strconv.Atoi(svc.port)
	if err != nil {
		return errors.New("端口格式错误")
	}
	if err := utils.IsUse(svc.bind, port); err != nil {
		return err
	}
	if common.State.Get(svc.key) {
		return errors.New("P2P隧道已在运行中")
	}
	if err := svc.run(); err != nil {
		return err
	}
	common.State.Set(svc.key, true)
	return
}

func (svc *P2P) Stop() {
	common.State.Set(svc.key, false)
	registry.Del(svc.key)
}

func (svc *P2P) IsRunning() bool {
	return common.State.Get(svc.key)
}

func (svc *P2P) run() (err error) {
	if svc.key == "" {
		return errors.New("please entry key")
	}
	data, err := common.GetP2PTunnelConfig(svc.httpUrl + "/api/v1/public/p2p/visit?key=" + svc.key)
	if err != nil {
		common.Logger.AddLog("p2p", svc.key+"获取配置失败，err:"+err.Error())
		return err
	}
	data.XTCPCfg.BindAddr = svc.bind
	data.XTCPCfg.BindPort, _ = strconv.Atoi(svc.port)

	// 是否禁用中继转发
	if data.DisableForward == 1 {
		data.XTCPCfg.FallbackTo = ""
	}

	registry.Del(svc.key)
	s := frpc.NewService(data.Common, nil, []v1.VisitorConfigurer{
		&data.XTCPCfg,
		&data.STCPCfg,
	})
	if err := s.Start(); err != nil {
		common.Logger.AddLog("p2p", svc.key+"启动失败，err:"+err.Error())
		return err
	}
	_ = registry.Set(svc.key, s)
	return nil
}
