package task

import (
	"server/model"
	"server/repository"
	"server/repository/cache"
	"server/service/common/node_port"
	"server/service/engine"
	"time"
)

func tunnelClean() {
	// 过期超过30天的隧道，自动删除
	now := time.Now().Add(-30 * 24 * time.Hour).Unix()
	db, _, _ := repository.Get("")
	hosts, _ := db.GostClientHost.Preload(db.GostClientHost.Node).Where(
		db.GostClientHost.ChargingType.Eq(model.GOST_CONFIG_CHARGING_CUCLE_DAY),
		db.GostClientHost.ExpAt.Lte(now),
	).Find()
	for _, host := range hosts {
		if host.CustomDomain != "" {
			_, _ = db.GostClientHostDomain.Where(db.GostClientHostDomain.Domain.Eq(host.CustomDomain)).Delete()
		}
		_, _ = db.GostNodeDomain.Where(
			db.GostNodeDomain.Prefix.Eq(host.DomainPrefix),
			db.GostNodeDomain.NodeCode.Eq(host.NodeCode),
		).Delete()
		_, _ = db.GostAuth.Where(db.GostAuth.TunnelCode.Eq(host.Code)).Delete()
		_, _ = db.GostClientHost.Where(db.GostClientHost.Code.Eq(host.Code)).Delete()

		engine.ClientRemoveHostConfig(db, *host, host.Node)
		cache.DelTunnelInfo(host.Code)
		cache.DelAdmissionInfo(host.Code)
	}

	forwards, _ := db.GostClientForward.Preload(db.GostClientForward.Node).Where(
		db.GostClientForward.ChargingType.Eq(model.GOST_CONFIG_CHARGING_CUCLE_DAY),
		db.GostClientForward.ExpAt.Lte(now),
	).Find()
	for _, forward := range forwards {
		_, _ = db.GostNodePort.Where(db.GostNodePort.Port.Eq(forward.Port), db.GostNodePort.NodeCode.Eq(forward.NodeCode)).Delete()
		_, _ = db.GostAuth.Where(db.GostAuth.TunnelCode.Eq(forward.Code)).Delete()
		_, _ = db.GostClientForward.Where(db.GostClientForward.Code.Eq(forward.Code)).Delete()
		node_port.ReleasePort(forward.NodeCode, forward.Port)
		engine.ClientRemoveForwardConfig(*forward)
		cache.DelTunnelInfo(forward.Code)
		cache.DelAdmissionInfo(forward.Code)
	}

	tunnels, _ := db.GostClientTunnel.Preload(db.GostClientTunnel.Node).Where(
		db.GostClientTunnel.ChargingType.Eq(model.GOST_CONFIG_CHARGING_CUCLE_DAY),
		db.GostClientTunnel.ExpAt.Lte(now),
	).Find()
	for _, tunnel := range tunnels {
		_, _ = db.GostAuth.Where(db.GostAuth.TunnelCode.Eq(tunnel.Code)).Delete()
		_, _ = db.GostClientTunnel.Where(db.GostClientTunnel.Code.Eq(tunnel.Code)).Delete()
		engine.ClientRemoveTunnelConfig(db, *tunnel, tunnel.Node)
		cache.DelTunnelInfo(tunnel.Code)
	}

	p2ps, _ := db.GostClientP2P.Preload(db.GostClientP2P.Node).Where(
		db.GostClientP2P.ChargingType.Eq(model.GOST_CONFIG_CHARGING_CUCLE_DAY),
		db.GostClientP2P.ExpAt.Lte(now),
	).Find()
	for _, p2p := range p2ps {
		_, _ = db.GostAuth.Where(db.GostAuth.TunnelCode.Eq(p2p.Code)).Delete()
		_, _ = db.GostClientP2P.Where(db.GostClientP2P.Code.Eq(p2p.Code)).Delete()
		engine.ClientRemoveP2PConfig(*p2p)
		cache.DelTunnelInfo(p2p.Code)
	}

	proxys, _ := db.GostClientProxy.Preload(db.GostClientProxy.Node).Where(
		db.GostClientProxy.ChargingType.Eq(model.GOST_CONFIG_CHARGING_CUCLE_DAY),
		db.GostClientProxy.ExpAt.Lte(now),
	).Find()
	for _, proxy := range proxys {
		_, _ = db.GostAuth.Where(db.GostAuth.TunnelCode.Eq(proxy.Code)).Delete()
		_, _ = db.GostClientProxy.Where(db.GostClientProxy.Code.Eq(proxy.Code)).Delete()
		engine.ClientRemoveProxyConfig(*proxy)
		cache.DelTunnelInfo(proxy.Code)
	}
}
