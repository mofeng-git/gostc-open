package todo

import (
	"server/bootstrap"
	"server/model"
	"server/repository"
	"server/service/common/cache"
)

func init() {
	bootstrap.AddTodo(func() {
		db, _, _ := repository.Get("")
		var authList []model.GostAuth
		db.Find(&authList)
		var authMap = make(map[string]model.GostAuth)
		for _, item := range authList {
			authMap[item.TunnelCode] = item
		}

		var hosts []model.GostClientHost
		db.Find(&hosts)
		for _, host := range hosts {
			cache.SetTunnelInfo(cache.TunnelInfo{
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
			cache.SetGostAuth(auth.User, auth.Password, host.Code)
		}

		var forwards []model.GostClientForward
		db.Find(&forwards)
		for _, forward := range forwards {
			cache.SetTunnelInfo(cache.TunnelInfo{
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
			cache.SetGostAuth(auth.User, auth.Password, forward.Code)
		}

		var tunnels []model.GostClientTunnel
		db.Find(&tunnels)
		for _, tunnel := range tunnels {
			cache.SetTunnelInfo(cache.TunnelInfo{
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
			cache.SetGostAuth(auth.User, auth.Password, tunnel.Code)
		}
	})

	//bootstrap.AddTodo(func() {
	//	db, _, _ := repository.Get("")
	//	{
	//		var codes []string
	//		db.Model(&model.GostClientHost{}).Pluck("code", &codes)
	//		db.Model(&model.GostAuth{}).Where("tunnel_code in ?", codes).Update("tunnel_type", model.GOST_TUNNEL_TYPE_HOST)
	//	}
	//	{
	//		var codes []string
	//		db.Model(&model.GostClientForward{}).Pluck("code", &codes)
	//		db.Model(&model.GostAuth{}).Where("tunnel_code in ?", codes).Update("tunnel_type", model.GOST_TUNNEL_TYPE_FORWARD)
	//	}
	//	{
	//		var codes []string
	//		db.Model(&model.GostClientTunnel{}).Pluck("code", &codes)
	//		db.Model(&model.GostAuth{}).Where("tunnel_code in ?", codes).Update("tunnel_type", model.GOST_TUNNEL_TYPE_TUNNEL)
	//	}
	//})
}
