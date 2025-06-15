package todo

import (
	cache2 "github.com/patrickmn/go-cache"
	"server/model"
	"server/repository"
	cache3 "server/repository/cache"
)

func gostClient() {
	db, _, _ := repository.Get("")
	var clientCodes []string
	_ = db.GostClient.Pluck(db.GostClient.Code, &clientCodes)
	for _, code := range clientCodes {
		cache3.SetClientOnline(code, false, cache2.NoExpiration)
	}

	authList, _ := db.GostAuth.Find()
	var authMap = make(map[string]model.GostAuth)
	for _, item := range authList {
		authMap[item.TunnelCode] = *item
	}

	hosts, _ := db.GostClientHost.Find()
	for _, host := range hosts {
		cache3.SetTunnelInfo(cache3.TunnelInfo{
			Code:        host.Code,
			Type:        model.GOST_TUNNEL_TYPE_HOST,
			ClientCode:  host.ClientCode,
			UserCode:    host.UserCode,
			NodeCode:    host.NodeCode,
			ChargingTye: host.ChargingType,
			ExpAt:       host.ExpAt,
			Limiter:     host.Limiter,
		})
		auth := authMap[host.Code]
		cache3.SetGostAuth(auth.User, auth.Password, host.Code)
	}

	forwards, _ := db.GostClientForward.Find()
	for _, forward := range forwards {
		cache3.SetTunnelInfo(cache3.TunnelInfo{
			Code:        forward.Code,
			Type:        model.GOST_TUNNEL_TYPE_FORWARD,
			ClientCode:  forward.ClientCode,
			UserCode:    forward.UserCode,
			NodeCode:    forward.NodeCode,
			ChargingTye: forward.ChargingType,
			ExpAt:       forward.ExpAt,
			Limiter:     forward.Limiter,
		})
		auth := authMap[forward.Code]
		cache3.SetGostAuth(auth.User, auth.Password, forward.Code)
	}

	tunnels, _ := db.GostClientTunnel.Find()
	for _, tunnel := range tunnels {
		cache3.SetTunnelInfo(cache3.TunnelInfo{
			Code:        tunnel.Code,
			Type:        model.GOST_TUNNEL_TYPE_TUNNEL,
			ClientCode:  tunnel.ClientCode,
			UserCode:    tunnel.UserCode,
			NodeCode:    tunnel.NodeCode,
			ChargingTye: tunnel.ChargingType,
			ExpAt:       tunnel.ExpAt,
			Limiter:     tunnel.Limiter,
		})
		auth := authMap[tunnel.Code]
		cache3.SetGostAuth(auth.User, auth.Password, tunnel.Code)
	}

	proxys, _ := db.GostClientProxy.Find()
	for _, proxy := range proxys {
		cache3.SetTunnelInfo(cache3.TunnelInfo{
			Code:        proxy.Code,
			Type:        model.GOST_TUNNEL_TYPE_PROXY,
			ClientCode:  proxy.ClientCode,
			UserCode:    proxy.UserCode,
			NodeCode:    proxy.NodeCode,
			ChargingTye: proxy.ChargingType,
			ExpAt:       proxy.ExpAt,
			Limiter:     proxy.Limiter,
		})
		auth := authMap[proxy.Code]
		cache3.SetGostAuth(auth.User, auth.Password, proxy.Code)
	}

	p2ps, _ := db.GostClientP2P.Find()
	for _, p2p := range p2ps {
		cache3.SetTunnelInfo(cache3.TunnelInfo{
			Code:        p2p.Code,
			Type:        model.GOST_TUNNEL_TYPE_P2P,
			ClientCode:  p2p.ClientCode,
			UserCode:    p2p.UserCode,
			NodeCode:    p2p.NodeCode,
			ChargingTye: p2p.ChargingType,
			ExpAt:       p2p.ExpAt,
			Limiter:     p2p.Limiter,
		})
		auth := authMap[p2p.Code]
		cache3.SetGostAuth(auth.User, auth.Password, p2p.Code)
	}
}
