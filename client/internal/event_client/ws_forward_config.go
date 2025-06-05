package event_client

import (
	"github.com/go-gost/core/logger"
	"github.com/go-gost/x/config"
	"github.com/go-gost/x/config/parsing/admission"
	"github.com/go-gost/x/config/parsing/chain"
	"github.com/go-gost/x/config/parsing/limiter"
	"github.com/go-gost/x/config/parsing/observer"
	"github.com/go-gost/x/config/parsing/service"
	"github.com/go-gost/x/registry"
)

type ClientForwardConfigData struct {
	SvcList       []config.ServiceConfig
	Chain         config.ChainConfig
	Limiter       config.LimiterConfig
	CLimiter      config.LimiterConfig
	RLimiter      config.LimiterConfig
	Obs           config.ObserverConfig
	AdmissionList []config.AdmissionConfig
}

func (e *Event) WsForwardConfig(data ClientForwardConfigData) {
	parseChain, err := chain.ParseChain(&data.Chain, logger.Default())
	if err == nil {
		registry.ChainRegistry().Unregister(data.Chain.Name)
		_ = registry.ChainRegistry().Register(data.Chain.Name, parseChain)
	}
	if data.Limiter.Name != "" {
		trafficLimiter := limiter.ParseTrafficLimiter(&data.Limiter)
		registry.TrafficLimiterRegistry().Unregister(data.Limiter.Name)
		_ = registry.TrafficLimiterRegistry().Register(data.Limiter.Name, trafficLimiter)
	}
	if data.CLimiter.Name != "" {
		connLimiter := limiter.ParseConnLimiter(&data.CLimiter)
		registry.ConnLimiterRegistry().Unregister(data.CLimiter.Name)
		_ = registry.ConnLimiterRegistry().Register(data.CLimiter.Name, connLimiter)
	}
	if data.RLimiter.Name != "" {
		rateLimiter := limiter.ParseRateLimiter(&data.RLimiter)
		registry.RateLimiterRegistry().Unregister(data.RLimiter.Name)
		_ = registry.RateLimiterRegistry().Register(data.RLimiter.Name, rateLimiter)
	}
	if data.Obs.Name != "" {
		parseObserver := observer.ParseObserver(&data.Obs)
		registry.ObserverRegistry().Unregister(data.Obs.Name)
		_ = registry.ObserverRegistry().Register(data.Obs.Name, parseObserver)
	}

	for _, item := range data.AdmissionList {
		if item.Name != "" {
			parseAdmission := admission.ParseAdmission(&item)
			registry.AdmissionRegistry().Unregister(item.Name)
			_ = registry.AdmissionRegistry().Register(item.Name, parseAdmission)
		}
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
		e.svcMap.Store(svcCfg.Name, true)
	}
}
