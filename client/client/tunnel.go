package main

import (
	"fmt"
	"github.com/go-gost/core/logger"
	"github.com/go-gost/x/config/parsing/chain"
	"github.com/go-gost/x/config/parsing/limiter"
	"github.com/go-gost/x/config/parsing/service"
	"github.com/go-gost/x/registry"
	"gostc-sub/common"
	"gostc-sub/pkg/signal"
	"strings"
)

type Tunnel struct {
	VKey string
	Port string
}

func runVisit(vTunnels string, apiurl string) {
	var tunnelList []Tunnel
	tunnels := strings.Split(vTunnels, ",")
	for _, tunnel := range tunnels {
		kp := strings.Split(tunnel, ":")
		if len(kp) != 2 {
			continue
		}
		tunnelList = append(tunnelList, Tunnel{
			VKey: kp[0],
			Port: kp[1],
		})
	}

	for _, tunnel := range tunnelList {
		data, err := common.GetVisitConfig(apiurl + "/api/v1/public/gost/visit?key=" + tunnel.VKey)
		if err != nil {
			fmt.Println("获取隧道配置失败", tunnel.VKey, tunnel.Port, err)
			continue
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
			svcCfg.Addr = ":" + tunnel.Port
			parseService, _ := service.ParseService(&svcCfg)
			go parseService.Serve()
			_ = registry.ServiceRegistry().Register(svcCfg.Name, parseService)
		}
		fmt.Println("隧道启动成功", tunnel.VKey, "0.0.0.0:"+tunnel.Port)
	}
	<-signal.Free()
}
