package warn_msg

import (
	"server/model"
	"time"
)

func GetHostWarnMsg(host model.GostClientHost) string {
	if host.Enable != 1 {
		return "已停用"
	}
	if host.ChargingType == model.GOST_CONFIG_CHARGING_CUCLE_DAY && host.ExpAt < time.Now().Unix() {
		return "已到期"
	}
	if host.Status == 2 {
		return "因未知原因被禁用"
	}
	return ""
}
func GetForwardWarnMsg(forward model.GostClientForward) string {
	if forward.Enable != 1 {
		return "已停用"
	}
	if forward.ChargingType == model.GOST_CONFIG_CHARGING_CUCLE_DAY && forward.ExpAt < time.Now().Unix() {
		return "已到期"
	}
	if forward.Status == 2 {
		return "因未知原因被禁用"
	}
	return ""
}

func GetTunnelWarnMsg(tunnel model.GostClientTunnel) string {
	if tunnel.Enable != 1 {
		return "已停用"
	}
	if tunnel.ChargingType == model.GOST_CONFIG_CHARGING_CUCLE_DAY && tunnel.ExpAt < time.Now().Unix() {
		return "已到期"
	}
	if tunnel.Status == 2 {
		return "因未知原因被禁用"
	}
	return ""
}

func GetProxyWarnMsg(proxy model.GostClientProxy) string {
	if proxy.Enable != 1 {
		return "已停用"
	}
	if proxy.ChargingType == model.GOST_CONFIG_CHARGING_CUCLE_DAY && proxy.ExpAt < time.Now().Unix() {
		return "已到期"
	}
	if proxy.Status == 2 {
		return "因未知原因被禁用"
	}
	return ""
}

func GetP2PWarnMsg(p2p model.GostClientP2P) string {
	if p2p.Enable != 1 {
		return "已停用"
	}
	if p2p.ChargingType == model.GOST_CONFIG_CHARGING_CUCLE_DAY && p2p.ExpAt < time.Now().Unix() {
		return "已到期"
	}
	if p2p.Status == 2 {
		return "因未知原因被禁用"
	}
	return ""
}
