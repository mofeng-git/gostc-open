package gost_engine

import (
	"github.com/go-gost/x/config"
	"github.com/lesismal/arpc"
	"server/model"
	"server/repository/query"
	"server/service/common/cache"
	"server/service/common/warn_msg"
	"time"
)

type ARpcClientEngine struct {
	code   string
	client *arpc.Client
}

func NewARpcClientEngine(code string, client *arpc.Client) *ARpcClientEngine {
	return &ARpcClientEngine{code: code, client: client}
}

func (e *ARpcClientEngine) Stop(msg string) {
	_ = e.client.Notify("stop", msg, time.Second*5)
}

func (e *ARpcClientEngine) PortCheck(tx *query.Query, ip, port string) (async bool, allowUse bool) {
	client, err := tx.GostClient.Where(tx.GostClient.Code.Eq(e.code)).First()
	if err != nil {
		return
	}
	var baseConfig model.SystemConfigBase
	cache.GetSystemConfigBase(&baseConfig)
	data := client.GenerateClientPortCheck(baseConfig.BaseUrl, port)
	var relay string
	_ = e.client.Call("port_check", data, &relay, time.Second*5)
	if relay == "success" {
		return true, true
	} else {
		return true, false
	}
}

func (e *ARpcClientEngine) HostConfig(tx *query.Query, hostCode string) {
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
	cache.SetIngress(host.DomainPrefix+"."+host.Node.Domain, host.Code)
	cache.SetIngress(host.DomainPrefix+"."+host.Node.Domain+":"+host.Node.TunnelInPort, host.Code)
	cache.SetIngress(host.DomainPrefix+"."+host.Node.Domain+":80", host.Code)
	cache.SetIngress(host.DomainPrefix+"."+host.Node.Domain+":443", host.Code)
	_ = e.client.Notify("host_config", data, time.Second*5)
	NodeIngress(tx, host.NodeCode)
}

func (e *ARpcClientEngine) RemoveHost(tx *query.Query, host model.GostClientHost, node model.GostNode) {
	cache.DelIngress(host.DomainPrefix + "." + node.Domain)
	cache.DelIngress(host.DomainPrefix + "." + node.Domain + ":" + node.TunnelInPort)
	cache.DelIngress(host.DomainPrefix + "." + node.Domain + ":80")
	cache.DelIngress(host.DomainPrefix + "." + node.Domain + ":443")
	_ = e.client.Notify("remove_config", []string{host.Code}, time.Second*5)
	NodeIngress(tx, node.Code)
}

func (e *ARpcClientEngine) ForwardConfig(tx *query.Query, forwardCode string) {
	forward, _ := tx.GostClientForward.Preload(tx.GostClientForward.Node).Where(tx.GostClientForward.Code.Eq(forwardCode)).First()
	if forward == nil {
		return
	}
	if warn_msg.GetForwardWarnMsg(*forward) != "" {
		ClientRemoveForwardConfig(*forward)
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
	_ = e.client.Notify("forward_config", data, time.Second*5)
}

func (e *ARpcClientEngine) RemoveForward(tx *query.Query, forward model.GostClientForward) {
	_ = e.client.Notify("remove_config", []string{
		"tcp_" + forward.Code,
		"udp_" + forward.Code,
	}, time.Second*5)
}

func (e *ARpcClientEngine) TunnelConfig(tx *query.Query, tunnelCode string) {
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
	cache.SetIngress(tunnel.Code, tunnel.Code)
	_ = e.client.Notify("tunnel_config", data, time.Second*5)
	NodeIngress(tx, tunnel.NodeCode)
}

func (e *ARpcClientEngine) RemoveTunnel(tx *query.Query, tunnel model.GostClientTunnel, node model.GostNode) {
	cache.DelIngress(tunnel.Code)
	_ = e.client.Notify("remove_config", []string{
		"tcp_" + tunnel.Code,
		"udp_" + tunnel.Code,
	}, time.Second*5)
	NodeIngress(tx, node.Code)
}

func (e *ARpcClientEngine) P2PConfig(tx *query.Query, p2pCode string) {
	p2p, _ := tx.GostClientP2P.Preload(tx.GostClientP2P.Node).Where(tx.GostClientP2P.Code.Eq(p2pCode)).First()
	if p2p == nil {
		return
	}
	if warn_msg.GetP2PWarnMsg(*p2p) != "" {
		ClientRemoveP2PConfig(*p2p)
		return
	}

	var data ClientP2PConfigData
	if p2p.Node.P2P == 1 {
		data.Code = p2pCode
		data.Common, _ = p2p.GenerateCommonCfg()
		data.STCPCfg, data.XTCPCfg = p2p.GenerateProxyCfgs()
	}
	_ = e.client.Notify("p2p_config", data, time.Second*5)
}

func (e *ARpcClientEngine) RemoveP2P(tx *query.Query, p2p model.GostClientP2P) {
	_ = e.client.Notify("remove_config", []string{
		p2p.Code,
	}, time.Second*5)
}

func (e *ARpcClientEngine) ProxyConfig(tx *query.Query, proxyCode string) {
	proxy, _ := tx.GostClientProxy.Preload(tx.GostClientProxy.Node).Where(tx.GostClientProxy.Code.Eq(proxyCode)).First()
	if proxy == nil {
		return
	}
	if warn_msg.GetProxyWarnMsg(*proxy) != "" {
		ClientRemoveProxyConfig(*proxy)
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
	_ = e.client.Notify("proxy_config", data, time.Second*5)
}

func (e *ARpcClientEngine) RemoveProxy(tx *query.Query, proxy model.GostClientProxy) {
	_ = e.client.Notify("remove_config", []string{
		"proxy_" + proxy.Code,
	}, time.Second*5)
}
