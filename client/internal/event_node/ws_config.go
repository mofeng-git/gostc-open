package event_node

import (
	"github.com/go-gost/x/config"
	"github.com/go-gost/x/config/parsing/auth"
	"github.com/go-gost/x/config/parsing/ingress"
	"github.com/go-gost/x/config/parsing/limiter"
	"github.com/go-gost/x/config/parsing/observer"
	"github.com/go-gost/x/config/parsing/service"
	"github.com/go-gost/x/registry"
	"gostc-sub/pkg/p2p/frps"
	v1 "gostc-sub/pkg/p2p/pkg/config/v1"
	registry2 "gostc-sub/pkg/p2p/registry"
)

type ConfigData struct {
	SvcList    []config.ServiceConfig
	Auther     config.AutherConfig
	Ingress    config.IngressConfig
	Limiter    config.LimiterConfig
	Obs        config.ObserverConfig
	P2PCfgCode string
	P2PCfg     v1.ServerConfig
}

func (e *Event) WsConfig(data ConfigData) {
	if data.Limiter.Name != "" {
		trafficLimiter := limiter.ParseTrafficLimiter(&data.Limiter)
		registry.TrafficLimiterRegistry().Unregister(data.Limiter.Name)
		_ = registry.TrafficLimiterRegistry().Register(data.Limiter.Name, trafficLimiter)
	}
	if data.Ingress.Name != "" {
		parseIngress := ingress.ParseIngress(&data.Ingress)
		registry.IngressRegistry().Unregister(data.Ingress.Name)
		_ = registry.IngressRegistry().Register(data.Ingress.Name, parseIngress)
	}
	if data.Auther.Name != "" {
		parseAuther := auth.ParseAuther(&data.Auther)
		registry.AutherRegistry().Unregister(data.Auther.Name)
		_ = registry.AutherRegistry().Register(data.Auther.Name, parseAuther)
	}
	if data.Obs.Name != "" {
		parseObserver := observer.ParseObserver(&data.Obs)
		registry.ObserverRegistry().Unregister(data.Obs.Name)
		_ = registry.ObserverRegistry().Register(data.Obs.Name, parseObserver)
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

	if data.P2PCfgCode != "" {
		registry2.Del(data.P2PCfgCode)
		svc := frps.NewService(data.P2PCfg)
		if err := svc.Start(); err == nil {
			_ = registry2.Set(data.P2PCfgCode, svc)
		}
	}
}
