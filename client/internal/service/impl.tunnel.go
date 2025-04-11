package service

import (
	"errors"
	"github.com/go-gost/core/logger"
	"github.com/go-gost/x/config/parsing/chain"
	"github.com/go-gost/x/config/parsing/limiter"
	"github.com/go-gost/x/config/parsing/service"
	"github.com/go-gost/x/registry"
	"gostc-sub/internal/common"
	"gostc-sub/pkg/utils"
	"strconv"
)

type Tunnel struct {
	key      string
	httpUrl  string
	port     string
	svcNames []string
}

func NewTunnel(httpUrl, key, port string) *Tunnel {
	return &Tunnel{
		key:     key,
		httpUrl: httpUrl,
		port:    port,
	}
}

func (svc *Tunnel) Start() (err error) {
	if !utils.ValidatePort(svc.port) {
		return errors.New("本地端口格式错误")
	}
	port, err := strconv.Atoi(svc.port)
	if err != nil {
		return errors.New("端口格式错误")
	}
	if utils.IsUse(port) {
		return errors.New("本地端口已被占用")
	}
	if common.State.Get(svc.key) {
		return errors.New("私有隧道已在运行中")
	}
	if err := svc.run(); err != nil {
		return err
	}
	common.State.Set(svc.key, true)
	return
}

func (svc *Tunnel) Stop() {
	common.State.Set(svc.key, false)
	for _, name := range svc.svcNames {
		if s := registry.ServiceRegistry().Get(name); s != nil {
			_ = s.Close()
		}
		registry.ServiceRegistry().Unregister(name)
	}
}

func (svc *Tunnel) IsRunning() bool {
	return common.State.Get(svc.key)
}

func (svc *Tunnel) run() (err error) {
	if svc.key == "" {
		return errors.New("please entry key")
	}
	data, err := common.GetVisitTunnelConfig(svc.httpUrl + "/api/v1/public/gost/visit?key=" + svc.key)
	if err != nil {
		common.Logger.AddLog("tunnel", svc.key+"获取配置失败，err:"+err.Error())
		return err
	}
	parseChain, _ := chain.ParseChain(&data.Chain, logger.Default())
	_ = registry.ChainRegistry().Register(data.Chain.Name, parseChain)
	trafficLimiter := limiter.ParseTrafficLimiter(&data.Limiter)
	_ = registry.TrafficLimiterRegistry().Register(data.Limiter.Name, trafficLimiter)
	connLimiter := limiter.ParseConnLimiter(&data.CLimiter)
	_ = registry.ConnLimiterRegistry().Register(data.CLimiter.Name, connLimiter)
	rateLimiter := limiter.ParseRateLimiter(&data.RLimiter)
	_ = registry.RateLimiterRegistry().Register(data.RLimiter.Name, rateLimiter)
	for _, svcCfg := range data.SvcList {
		svcCfg.Addr = ":" + svc.port
		parseService, _ := service.ParseService(&svcCfg)
		svc.svcNames = append(svc.svcNames, svcCfg.Name)
		go parseService.Serve()
		_ = registry.ServiceRegistry().Register(svcCfg.Name, parseService)
	}
	return nil
}
