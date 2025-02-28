package gost_engine

import (
	"github.com/go-gost/x/config"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"server/model"
	"server/service/common/cache"
	"server/service/common/warn_msg"
)

func ClientStop(code string, msg string) {
	WriteMessage(code, NewMessage(uuid.NewString(), "stop", map[string]string{
		"reason": msg,
	}))
}

type ClientHostConfigData struct {
	Svc           config.ServiceConfig
	Chain         config.ChainConfig
	AdmissionList []config.AdmissionConfig
}

func ClientHostConfig(tx *gorm.DB, hostCode string) {
	var host model.GostClientHost
	if tx.Preload("Node").Where("code = ?", hostCode).First(&host).RowsAffected == 0 {
		return
	}
	if warn_msg.GetHostWarnMsg(host) != "" {
		ClientRemoveHostConfig(tx, host, host.Node)
		return
	}
	var auth model.GostAuth
	if tx.Where("tunnel_code = ?", hostCode).First(&auth).RowsAffected == 0 {
		return
	}
	var data ClientHostConfigData
	chain := host.GenerateChainConfig(auth)
	admissionWhite, admissionBlack := host.GenerateWhiteAdmission(), host.GenerateBlackAdmission()
	svcCfg, ok := host.GenerateSvcConfig(chain.Name, admissionWhite.Name, admissionBlack.Name)
	if !ok {
		return
	}
	data.Chain = chain
	data.Svc = svcCfg
	data.AdmissionList = []config.AdmissionConfig{admissionWhite, admissionBlack}
	WriteMessage(host.ClientCode, NewMessage(uuid.NewString(), "host_config", data))
	cache.SetIngress(host.DomainPrefix+"."+host.Node.Domain, host.Code)
	cache.SetIngress(host.DomainPrefix+"."+host.Node.Domain+":"+host.Node.TunnelInPort, host.Code)
	cache.SetIngress(host.DomainPrefix+"."+host.Node.Domain+":80", host.Code)
	cache.SetIngress(host.DomainPrefix+"."+host.Node.Domain+":443", host.Code)
	NodeIngress(tx, host.NodeCode)
}

func ClientRemoveHostConfig(tx *gorm.DB, host model.GostClientHost, node model.GostNode) {
	cache.DelIngress(host.DomainPrefix + "." + node.Domain)
	cache.DelIngress(host.DomainPrefix + "." + node.Domain + ":" + node.TunnelInPort)
	cache.DelIngress(host.DomainPrefix + "." + node.Domain + ":80")
	cache.DelIngress(host.DomainPrefix + "." + node.Domain + ":443")
	WriteMessage(host.ClientCode, NewMessage(uuid.NewString(), "remove_config", []string{host.Code}))
	NodeIngress(tx, node.Code)
}

type ClientForwardConfigData struct {
	SvcList       []config.ServiceConfig
	Chain         config.ChainConfig
	Limiter       config.LimiterConfig
	CLimiter      config.LimiterConfig
	RLimiter      config.LimiterConfig
	Obs           config.ObserverConfig
	AdmissionList []config.AdmissionConfig
}

func ClientForwardConfig(tx *gorm.DB, forwardCode string) {
	var forward model.GostClientForward
	if tx.Preload("Node").Where("code = ?", forwardCode).First(&forward).RowsAffected == 0 {
		return
	}
	if warn_msg.GetForwardWarnMsg(forward) != "" {
		ClientRemoveForwardConfig(forward, forward.Node)
		return
	}
	var auth model.GostAuth
	if tx.Where("tunnel_code = ?", forwardCode).First(&auth).RowsAffected == 0 {
		return
	}

	var baseConfig model.SystemConfigBase
	cache.GetSystemConfigBase(&baseConfig)
	nodeVersion := cache.GetNodeVersion(forward.NodeCode)

	var data ClientForwardConfigData
	chain := forward.GenerateChainConfig(auth)
	limiter := forward.GenerateLimiter()
	rLimiter := forward.GenerateRLimiter()
	cLimiter := forward.GenerateCLimiter()
	obs := forward.GenerateObs(baseConfig.BaseUrl, nodeVersion)
	admissionWhite, admissionBlack := forward.GenerateWhiteAdmission(), forward.GenerateBlackAdmission()

	tcpSvcCfg, ok := forward.GenerateTcpSvcConfig(chain.Name, limiter.Name, cLimiter.Name, rLimiter.Name, obs.Name, admissionWhite.Name, admissionBlack.Name)
	if ok {
		data.SvcList = append(data.SvcList, tcpSvcCfg)
	}
	udpSvcCfg, ok := forward.GenerateUdpSvcConfig(chain.Name, limiter.Name, cLimiter.Name, rLimiter.Name, obs.Name, admissionWhite.Name, admissionBlack.Name)
	if ok {
		data.SvcList = append(data.SvcList, udpSvcCfg)
	}
	data.Chain = chain
	data.Limiter = limiter
	data.CLimiter = cLimiter
	data.RLimiter = rLimiter
	data.Obs = obs
	data.AdmissionList = []config.AdmissionConfig{admissionWhite, admissionBlack}
	WriteMessage(forward.ClientCode, NewMessage(uuid.NewString(), "forward_config", data))
}

func ClientRemoveForwardConfig(forward model.GostClientForward, node model.GostNode) {
	WriteMessage(forward.ClientCode, NewMessage(uuid.NewString(), "remove_config", []string{
		"tcp_" + forward.Code,
		"udp_" + forward.Code,
	}))
}

type ClientTunnelConfigData struct {
	SvcList []config.ServiceConfig
	Chain   config.ChainConfig
}

func ClientTunnelConfig(tx *gorm.DB, tunnelCode string) {
	var tunnel model.GostClientTunnel
	if tx.Preload("Node").Where("code = ?", tunnelCode).First(&tunnel).RowsAffected == 0 {
		return
	}
	if warn_msg.GetTunnelWarnMsg(tunnel) != "" {
		ClientRemoveTunnelConfig(tx, tunnel, tunnel.Node)
		return
	}
	var auth model.GostAuth
	if tx.Where("tunnel_code = ?", tunnelCode).First(&auth).RowsAffected == 0 {
		return
	}

	var data ClientTunnelConfigData
	chain := tunnel.GenerateChainConfig(auth)
	tcpSvcCfg, ok := tunnel.GenerateTcpSvcConfig(chain.Name)
	if ok {
		data.SvcList = append(data.SvcList, tcpSvcCfg)
	}
	udpSvcCfg, ok := tunnel.GenerateUdpSvcConfig(chain.Name)
	if ok {
		data.SvcList = append(data.SvcList, udpSvcCfg)
	}
	data.Chain = chain
	WriteMessage(tunnel.ClientCode, NewMessage(uuid.NewString(), "tunnel_config", data))
	cache.SetIngress(tunnel.Code, tunnel.Code)
	NodeIngress(tx, tunnel.NodeCode)
}

func ClientRemoveTunnelConfig(tx *gorm.DB, tunnel model.GostClientTunnel, node model.GostNode) {
	cache.DelIngress(tunnel.Code)
	WriteMessage(tunnel.Code, NewMessage(uuid.NewString(), "remove_config", []string{
		"tcp_" + tunnel.Code,
		"udp_" + tunnel.Code,
	}))
	NodeIngress(tx, node.Code)
}
