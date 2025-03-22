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

func NodeStop(code string, msg string) {
	WriteMessage(code, NewMessage(uuid.NewString(), "stop", map[string]string{
		"reason": msg,
	}))
}

type NodeConfigData struct {
	SvcList    []config.ServiceConfig
	Auther     config.AutherConfig
	Ingress    config.IngressConfig
	Limiter    config.LimiterConfig
	Obs        config.ObserverConfig
	P2PCfgCode string
	P2PCfg     v1.ServerConfig
}

func NodeIngress(tx *query.Query, code string) {
	node, _ := tx.GostNode.Where(tx.GostNode.Code.Eq(code)).First()
	if node == nil {
		return
	}

	hosts, _ := tx.GostClientHost.Where(tx.GostClientHost.NodeCode.Eq(code)).Find()
	var newHosts []*model.GostClientHost
	for _, host := range hosts {
		if warn_msg.GetHostWarnMsg(*host) != "" {
			continue
		}
		newHosts = append(newHosts, host)
	}

	tunnels, _ := tx.GostClientTunnel.Where(tx.GostClientTunnel.NodeCode.Eq(code)).Find()
	var newTunnels []*model.GostClientTunnel
	for _, tunnel := range tunnels {
		if warn_msg.GetTunnelWarnMsg(*tunnel) != "" {
			continue
		}
		newTunnels = append(newTunnels, tunnel)
	}
	var data NodeConfigData
	data.Ingress = node.GenerateIngress(newHosts, newTunnels)
	WriteMessage(code, NewMessage(uuid.NewString(), "config", data))
}

func NodeConfig(tx *query.Query, code string) {
	node, err := tx.GostNode.Where(tx.GostNode.Code.Eq(code)).First()
	if err != nil {
		return
	}

	var baseConfig model.SystemConfigBase
	cache.GetSystemConfigBase(&baseConfig)

	var data NodeConfigData
	auther := node.GenerateAuther(baseConfig.BaseUrl)
	hosts, _ := tx.GostClientHost.Where(tx.GostClientHost.NodeCode.Eq(node.Code)).Find()
	tunnels, _ := tx.GostClientTunnel.Where(tx.GostClientTunnel.NodeCode.Eq(node.Code)).Find()
	ingress := node.GenerateIngress(hosts, tunnels)
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
	WriteMessage(code, NewMessage(uuid.NewString(), "config", data))
}

func NodePortCheck(tx *query.Query, code string, port string) {
	node, err := tx.GostNode.Where(tx.GostNode.Code.Eq(code)).First()
	if err != nil {
		return
	}

	var baseConfig model.SystemConfigBase
	cache.GetSystemConfigBase(&baseConfig)
	data := node.GenerateNodePortCheck(baseConfig.BaseUrl, port)
	WriteMessage(code, NewMessage(uuid.NewString(), "port_check", data))
}
