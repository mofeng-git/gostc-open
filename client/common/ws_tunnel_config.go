package common

import (
	"github.com/go-gost/core/logger"
	"github.com/go-gost/x/config"
	"github.com/go-gost/x/config/parsing/chain"
	"github.com/go-gost/x/config/parsing/service"
	"github.com/go-gost/x/registry"
)

type ClientTunnelConfigData struct {
	SvcList []config.ServiceConfig
	Chain   config.ChainConfig
}

func WsTunnelConfig(data ClientTunnelConfigData) {
	parseChain, err := chain.ParseChain(&data.Chain, logger.Default())
	if err == nil {
		registry.ChainRegistry().Unregister(data.Chain.Name)
		_ = registry.ChainRegistry().Register(data.Chain.Name, parseChain)
	}
	for _, svcCfg := range data.SvcList {
		if oldSvc := registry.ServiceRegistry().Get(svcCfg.Name); oldSvc != nil {
			registry.ServiceRegistry().Unregister(svcCfg.Name)
			_ = oldSvc.Close()
		}
		svc, err := service.ParseService(&svcCfg)
		if err != nil {
			continue
		}
		go svc.Serve()
		if err = registry.ServiceRegistry().Register(svcCfg.Name, svc); err != nil {
			_ = svc.Close()
			continue
		}
		SvcMap[svcCfg.Name] = true
	}
}
