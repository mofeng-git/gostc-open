package engine

import (
	"github.com/go-gost/core/logger"
	"github.com/go-gost/x/config"
	"github.com/go-gost/x/config/parsing/admission"
	"github.com/go-gost/x/config/parsing/chain"
	"github.com/go-gost/x/config/parsing/service"
	"github.com/go-gost/x/registry"
)

type ClientHostConfigData struct {
	Svc           config.ServiceConfig
	Chain         config.ChainConfig
	AdmissionList []config.AdmissionConfig
}

func (e *Event) WsHostConfig(data ClientHostConfigData) {
	parseChain, err := chain.ParseChain(&data.Chain, logger.Default())
	if err == nil {
		registry.ChainRegistry().Unregister(data.Chain.Name)
		_ = registry.ChainRegistry().Register(data.Chain.Name, parseChain)
	}

	for _, item := range data.AdmissionList {
		if item.Name != "" {
			parseAdmission := admission.ParseAdmission(&item)
			registry.AdmissionRegistry().Unregister(item.Name)
			_ = registry.AdmissionRegistry().Register(item.Name, parseAdmission)
		}
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
		e.svcMap.Store(data.Svc.Name, true)
	}
}
