package gost_engine

import (
	"github.com/go-gost/x/config"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"server/global"
	"server/model"
	"server/service/common/cache"
	"server/service/common/warn_msg"
)

func NodeStop(code string, msg string) {
	WriteMessage(code, NewMessage(uuid.NewString(), "stop", map[string]string{
		"reason": msg,
	}))
}

type NodeConfigData struct {
	SvcList []config.ServiceConfig
	Auther  config.AutherConfig
	Ingress config.IngressConfig
	Limiter config.LimiterConfig
	Obs     config.ObserverConfig
}

func NodeIngress(tx *gorm.DB, code string) {
	var node model.GostNode
	if tx.Where("code = ?", code).First(&node).RowsAffected == 0 {
		return
	}

	var hosts []model.GostClientHost
	tx.Where("node_code = ?", node.Code).Find(&hosts)
	var newHosts []model.GostClientHost
	for _, host := range hosts {
		if warn_msg.GetHostWarnMsg(host) != "" {
			continue
		}
		newHosts = append(newHosts, host)
	}
	var tunnels []model.GostClientTunnel
	tx.Where("node_code = ?", node.Code).Find(&tunnels)
	var newTunnels []model.GostClientTunnel
	for _, tunnel := range tunnels {
		if warn_msg.GetTunnelWarnMsg(tunnel) != "" {
			continue
		}
		newTunnels = append(newTunnels, tunnel)
	}
	var data NodeConfigData
	data.Ingress = node.GenerateIngress(newHosts, newTunnels)
	WriteMessage(code, NewMessage(uuid.NewString(), "config", data))
}

func NodeConfig(tx *gorm.DB, code string) {
	var node model.GostNode
	if global.DB.GetDB().Where("code = ?", code).First(&node).RowsAffected == 0 {
		return
	}

	var baseConfig model.SystemConfigBase
	cache.GetSystemConfigBase(&baseConfig)

	var data NodeConfigData
	auther := node.GenerateAuther(baseConfig.BaseUrl)

	var hosts []model.GostClientHost
	tx.Where("node_code = ?", node.Code).Find(&hosts)
	var tunnels []model.GostClientTunnel
	tx.Where("node_code = ?", node.Code).Find(&tunnels)
	ingress := node.GenerateIngress(hosts, tunnels)

	limiter := node.GenerateLimiter(baseConfig.BaseUrl)
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
	WriteMessage(code, NewMessage(uuid.NewString(), "config", data))
}
