package gost_engine

import (
	"github.com/go-gost/x/config"
	"github.com/google/uuid"
	"server/model"
	v1 "server/pkg/p2p_cfg/v1"
	"server/repository/query"
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

func ClientHostConfig(tx *query.Query, hostCode string) {
	host, _ := tx.GostClientHost.Preload(tx.GostClientHost.Node).Where(tx.GostClientHost.Code.Eq(hostCode)).First()
	if host == nil {
		return
	}

	if warn_msg.GetHostWarnMsg(*host) != "" {
		ClientRemoveHostConfig(tx, *host, host.Node)
		return
	}

	auth, _ := tx.GostAuth.Where(tx.GostAuth.TunnelCode.Eq(hostCode)).First()
	if auth == nil {
		return
	}
	var data ClientHostConfigData
	chain := host.GenerateChainConfig(*auth)
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

func ClientRemoveHostConfig(tx *query.Query, host model.GostClientHost, node model.GostNode) {
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

func ClientForwardConfig(tx *query.Query, forwardCode string) {
	forward, _ := tx.GostClientForward.Preload(tx.GostClientForward.Node).Where(tx.GostClientForward.Code.Eq(forwardCode)).First()
	if forward == nil {
		return
	}
	if warn_msg.GetForwardWarnMsg(*forward) != "" {
		ClientRemoveForwardConfig(*forward, forward.Node)
		return
	}

	auth, _ := tx.GostAuth.Where(tx.GostAuth.TunnelCode.Eq(forwardCode)).First()
	if auth == nil {
		return
	}

	var baseConfig model.SystemConfigBase
	cache.GetSystemConfigBase(&baseConfig)
	nodeVersion := cache.GetNodeVersion(forward.NodeCode)

	var data ClientForwardConfigData
	chain := forward.GenerateChainConfig(*auth)
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

func ClientTunnelConfig(tx *query.Query, tunnelCode string) {
	tunnel, _ := tx.GostClientTunnel.Preload(tx.GostClientTunnel.Node).Where(tx.GostClientTunnel.Code.Eq(tunnelCode)).First()
	if tunnel == nil {
		return
	}
	if warn_msg.GetTunnelWarnMsg(*tunnel) != "" {
		ClientRemoveTunnelConfig(tx, *tunnel, tunnel.Node)
		return
	}

	auth, _ := tx.GostAuth.Where(tx.GostAuth.TunnelCode.Eq(tunnelCode)).First()
	if auth == nil {
		return
	}

	var data ClientTunnelConfigData
	chain := tunnel.GenerateChainConfig(*auth)
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

func ClientRemoveTunnelConfig(tx *query.Query, tunnel model.GostClientTunnel, node model.GostNode) {
	cache.DelIngress(tunnel.Code)
	WriteMessage(tunnel.Code, NewMessage(uuid.NewString(), "remove_config", []string{
		"tcp_" + tunnel.Code,
		"udp_" + tunnel.Code,
	}))
	NodeIngress(tx, node.Code)
}

type ClientProxyConfigData struct {
	Svc      config.ServiceConfig
	Chain    config.ChainConfig
	Limiter  config.LimiterConfig
	CLimiter config.LimiterConfig
	RLimiter config.LimiterConfig
	Obs      config.ObserverConfig
}

func ClientProxyConfig(tx *query.Query, proxyCode string) {
	proxy, _ := tx.GostClientProxy.Preload(tx.GostClientProxy.Node).Where(tx.GostClientProxy.Code.Eq(proxyCode)).First()
	if proxy == nil {
		return
	}
	if warn_msg.GetProxyWarnMsg(*proxy) != "" {
		ClientRemoveProxyConfig(*proxy, proxy.Node)
		return
	}

	auth, _ := tx.GostAuth.Where(tx.GostAuth.TunnelCode.Eq(proxyCode)).First()
	if auth == nil {
		return
	}

	var baseConfig model.SystemConfigBase
	cache.GetSystemConfigBase(&baseConfig)
	nodeVersion := cache.GetNodeVersion(proxy.NodeCode)

	var data ClientProxyConfigData
	chain := proxy.GenerateChainConfig(*auth)
	limiter := proxy.GenerateLimiter()
	rLimiter := proxy.GenerateRLimiter()
	cLimiter := proxy.GenerateCLimiter()
	obs := proxy.GenerateObs(baseConfig.BaseUrl, nodeVersion)

	svcCfg, ok := proxy.GenerateSvcConfig(chain.Name, limiter.Name, cLimiter.Name, rLimiter.Name, obs.Name)
	if ok {
		data.Svc = svcCfg
	}
	data.Chain = chain
	data.Limiter = limiter
	data.CLimiter = cLimiter
	data.RLimiter = rLimiter
	data.Obs = obs
	WriteMessage(proxy.ClientCode, NewMessage(uuid.NewString(), "proxy_config", data))
}

func ClientRemoveProxyConfig(proxy model.GostClientProxy, node model.GostNode) {
	WriteMessage(proxy.ClientCode, NewMessage(uuid.NewString(), "remove_config", []string{
		"proxy_" + proxy.Code,
	}))
}

type ClientP2PConfigData struct {
	Code    string
	Common  v1.ClientCommonConfig
	STCPCfg v1.STCPProxyConfig
	XTCPCfg v1.XTCPProxyConfig
}

func ClientP2PConfig(tx *query.Query, p2pCode string) {
	p2p, _ := tx.GostClientP2P.Preload(tx.GostClientP2P.Node).Where(tx.GostClientP2P.Code.Eq(p2pCode)).First()
	if p2p == nil {
		return
	}
	if warn_msg.GetP2PWarnMsg(*p2p) != "" {
		ClientRemoveP2PConfig(*p2p, p2p.Node)
		return
	}

	var data ClientP2PConfigData
	if p2p.Node.P2P == 1 {
		data.Code = p2pCode
		data.Common, _ = p2p.GenerateCommonCfg()
		data.STCPCfg, data.XTCPCfg = p2p.GenerateProxyCfgs()
	}
	WriteMessage(p2p.ClientCode, NewMessage(uuid.NewString(), "p2p_config", data))
}

func ClientRemoveP2PConfig(p2p model.GostClientP2P, node model.GostNode) {
	WriteMessage(p2p.ClientCode, NewMessage(uuid.NewString(), "remove_config", []string{
		p2p.Code,
	}))
}

func ClientPortCheck(tx *query.Query, code string, port string) {
	client, err := tx.GostClient.Where(tx.GostClient.Code.Eq(code)).First()
	if err != nil {
		return
	}

	var baseConfig model.SystemConfigBase
	cache.GetSystemConfigBase(&baseConfig)
	data := client.GenerateClientPortCheck(baseConfig.BaseUrl, port)
	WriteMessage(code, NewMessage(uuid.NewString(), "port_check", data))
}
