package common

import (
	"github.com/go-gost/core/logger"
	"github.com/go-gost/x/config"
	"github.com/go-gost/x/config/parsing/chain"
	"github.com/go-gost/x/config/parsing/service"
	"github.com/go-gost/x/registry"
)

type ClientProxyConfigData struct {
	Svc      config.ServiceConfig
	Chain    config.ChainConfig
	Limiter  config.LimiterConfig
	CLimiter config.LimiterConfig
	RLimiter config.LimiterConfig
	Obs      config.ObserverConfig
}

func WsProxyConfig(data ClientProxyConfigData) {
	parseChain, err := chain.ParseChain(&data.Chain, logger.Default())
	if err == nil {
		registry.ChainRegistry().Unregister(data.Chain.Name)
		_ = registry.ChainRegistry().Register(data.Chain.Name, parseChain)
	}
	if oldSvc := registry.ServiceRegistry().Get(data.Svc.Name); oldSvc != nil {
		registry.ServiceRegistry().Unregister(data.Svc.Name)
		_ = oldSvc.Close()
	}
	svc, err := service.ParseService(&data.Svc)
	if err == nil {
		go svc.Serve()
		if err = registry.ServiceRegistry().Register(data.Svc.Name, svc); err != nil {
			_ = svc.Close()
		}
		SvcMap[data.Svc.Name] = true
	}
}
