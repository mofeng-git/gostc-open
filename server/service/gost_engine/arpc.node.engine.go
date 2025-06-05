package gost_engine

import (
	"fmt"
	"github.com/lesismal/arpc"
	"server/model"
	"server/repository/query"
	"server/service/common/cache"
	"server/service/common/warn_msg"
	"time"
)

func NewARpcNodeEngine(code, ip string, client *arpc.Client) *ARpcNodeEngine {
	return &ARpcNodeEngine{
		code:   code,
		ip:     ip,
		client: client,
	}
}

type ARpcNodeEngine struct {
	code   string
	ip     string
	client *arpc.Client
}

func (e *ARpcNodeEngine) PortCheck(tx *query.Query, ip, port string) (async bool, allowUse bool) {
	node, err := tx.GostNode.Where(tx.GostNode.Code.Eq(e.code)).First()
	if err != nil {
		return
	}
	var baseConfig model.SystemConfigBase
	cache.GetSystemConfigBase(&baseConfig)
	data := node.GenerateNodePortCheck(baseConfig.BaseUrl, port)
	var relay string
	_ = e.client.Call("port_check", data, &relay, time.Second*5)
	if relay == "success" {
		return true, true
	} else {
		return true, false
	}
}

func (e *ARpcNodeEngine) Config(tx *query.Query) {
	node, err := tx.GostNode.Where(tx.GostNode.Code.Eq(e.code)).First()
	if err != nil {
		return
	}

	var baseConfig model.SystemConfigBase
	cache.GetSystemConfigBase(&baseConfig)

	var data NodeConfigData
	auther := node.GenerateAuther(baseConfig.BaseUrl)
	hosts, _ := tx.GostClientHost.Where(tx.GostClientHost.NodeCode.Eq(node.Code)).Find()
	tunnels, _ := tx.GostClientTunnel.Where(tx.GostClientTunnel.NodeCode.Eq(node.Code)).Find()
	ingress := node.GenerateIngress(hosts, tunnels, cache.GetNodeCustomDomain(e.code))
	limiter := node.GenerateLimiter(baseConfig.BaseUrl)
	p2pCfg := node.GenerateP2PServiceConfig(baseConfig.BaseUrl)
	obs := node.GenerateObs(baseConfig.BaseUrl)
	tunnelAndHostSvcCfg, ok := node.GenerateTunnelAndHostServiceConfig(limiter.Name, auther.Name, ingress.Name, obs.Name)
	if ok {
		data.SvcList = append(data.SvcList, tunnelAndHostSvcCfg)
	}
	forwardSvcCfg, ok := node.GenerateForwardServiceConfig(limiter.Name, auther.Name, obs.Name)
	if ok {
		data.SvcList = append(data.SvcList, forwardSvcCfg)
	}
	if len(data.SvcList) == 0 {
		return
	}
	data.Auther = auther
	data.Ingress = ingress
	data.Limiter = limiter
	data.Obs = obs
	if node.P2P == 1 {
		data.P2PCfgCode = node.Code
		data.P2PCfg = p2pCfg
	}
	_ = e.client.Notify("config", data, time.Second*5)
}

func (e *ARpcNodeEngine) Ingress(tx *query.Query) {
	node, _ := tx.GostNode.Where(tx.GostNode.Code.Eq(e.code)).First()
	if node == nil {
		return
	}

	hosts, _ := tx.GostClientHost.Where(tx.GostClientHost.NodeCode.Eq(e.code)).Find()
	var newHosts []*model.GostClientHost
	for _, host := range hosts {
		if warn_msg.GetHostWarnMsg(*host) != "" {
			continue
		}
		newHosts = append(newHosts, host)
	}

	tunnels, _ := tx.GostClientTunnel.Where(tx.GostClientTunnel.NodeCode.Eq(e.code)).Find()
	var newTunnels []*model.GostClientTunnel
	for _, tunnel := range tunnels {
		if warn_msg.GetTunnelWarnMsg(*tunnel) != "" {
			continue
		}
		newTunnels = append(newTunnels, tunnel)
	}
	var data NodeConfigData
	data.Ingress = node.GenerateIngress(newHosts, newTunnels, cache.GetNodeCustomDomain(e.code))
	_ = e.client.Notify("config", data, time.Second*5)
}

func (e *ARpcNodeEngine) CustomDomain(tx *query.Query, domain, cert, key string, forceHttps int) {
	node, err := tx.GostNode.Where(tx.GostNode.Code.Eq(e.code)).First()
	if err != nil {
		return
	}
	_ = e.client.Notify("https_domain", HttpsDomainData{
		Domain:     domain,
		Target:     fmt.Sprintf("http://127.0.0.1:%s", node.TunnelInPort),
		Cert:       cert,
		Key:        key,
		ForceHttps: forceHttps,
	}, time.Second*5)
}

func (e *ARpcNodeEngine) Stop(msg string) {
	_ = e.client.Notify("stop", msg, time.Second*5)
}
