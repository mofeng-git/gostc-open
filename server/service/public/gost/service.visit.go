package service

import (
	"github.com/go-gost/x/config"
	"server/model"
	"server/repository"
)

type VisitCfgResp struct {
	SvcList  []config.ServiceConfig
	Chain    config.ChainConfig
	Limiter  config.LimiterConfig
	RLimiter config.LimiterConfig
	CLimiter config.LimiterConfig
}

func (service *service) VisitCfg(key string) (result VisitCfgResp) {
	db, _, _ := repository.Get("")
	var tunnel model.GostClientTunnel
	if db.Preload("Node").Where("v_key = ?", key).First(&tunnel).RowsAffected == 0 {
		return result
	}
	var auth model.GostAuth
	if db.Where("tunnel_code = ?", tunnel.Code).First(&auth).RowsAffected == 0 {
		return result
	}

	limiter := tunnel.GenerateVisitLimiter()
	cLimiter := tunnel.GenerateVisitCLimiter()
	rLimiter := tunnel.GenerateVisitRLimiter()
	chain := tunnel.GenerateVisitChainConfig(auth)
	tcpSvcConfig, ok := tunnel.GenerateVisitTcpSvcConfig(chain.Name, limiter.Name, cLimiter.Name, rLimiter.Name)
	if ok {
		result.SvcList = append(result.SvcList, tcpSvcConfig)
	}
	udpSvcConfig, ok := tunnel.GenerateVisitUdpSvcConfig(chain.Name, limiter.Name, cLimiter.Name, rLimiter.Name)
	if ok {
		result.SvcList = append(result.SvcList, udpSvcConfig)
	}
	result.Limiter = limiter
	result.CLimiter = cLimiter
	result.RLimiter = rLimiter
	result.Chain = chain
	return result
}
