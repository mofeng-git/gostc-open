package engine

import (
	"server/model"
	"server/repository/query"
)

func ClientStop(code string, msg string) {
	engine, ok := EngineRegistry.Get(code)
	if !ok {
		return
	}
	engine.GetClient().Stop(msg)
}

func ClientHostConfig(tx *query.Query, hostCode string) {
	host, _ := tx.GostClientHost.Where(tx.GostClientHost.Code.Eq(hostCode)).First()
	if host == nil {
		return
	}
	engine, ok := EngineRegistry.Get(host.ClientCode)
	if !ok {
		return
	}
	engine.GetClient().HostConfig(tx, hostCode)
}

func ClientRemoveHostConfig(tx *query.Query, host model.GostClientHost, node model.GostNode) {
	engine, ok := EngineRegistry.Get(host.ClientCode)
	if !ok {
		return
	}
	engine.GetClient().RemoveHost(tx, host, node)
}

func ClientForwardConfig(tx *query.Query, forwardCode string) {
	forward, _ := tx.GostClientForward.Preload(tx.GostClientForward.Node).Where(tx.GostClientForward.Code.Eq(forwardCode)).First()
	if forward == nil {
		return
	}
	engine, ok := EngineRegistry.Get(forward.ClientCode)
	if !ok {
		return
	}
	engine.GetClient().ForwardConfig(tx, forwardCode)
}

func ClientRemoveForwardConfig(forward model.GostClientForward) {
	engine, ok := EngineRegistry.Get(forward.ClientCode)
	if !ok {
		return
	}
	engine.GetClient().RemoveForward(nil, forward)
}

func ClientTunnelConfig(tx *query.Query, tunnelCode string) {
	tunnel, _ := tx.GostClientTunnel.Preload(tx.GostClientTunnel.Node).Where(tx.GostClientTunnel.Code.Eq(tunnelCode)).First()
	if tunnel == nil {
		return
	}
	engine, ok := EngineRegistry.Get(tunnel.ClientCode)
	if !ok {
		return
	}
	engine.GetClient().TunnelConfig(tx, tunnelCode)
}

func ClientRemoveTunnelConfig(tx *query.Query, tunnel model.GostClientTunnel, node model.GostNode) {
	engine, ok := EngineRegistry.Get(tunnel.ClientCode)
	if !ok {
		return
	}
	engine.GetClient().RemoveTunnel(tx, tunnel, node)
}

func ClientProxyConfig(tx *query.Query, proxyCode string) {
	proxy, _ := tx.GostClientProxy.Preload(tx.GostClientProxy.Node).Where(tx.GostClientProxy.Code.Eq(proxyCode)).First()
	if proxy == nil {
		return
	}
	engine, ok := EngineRegistry.Get(proxy.ClientCode)
	if !ok {
		return
	}
	engine.GetClient().ProxyConfig(tx, proxyCode)
}

func ClientRemoveProxyConfig(proxy model.GostClientProxy) {
	engine, ok := EngineRegistry.Get(proxy.ClientCode)
	if !ok {
		return
	}
	engine.GetClient().RemoveProxy(nil, proxy)
}

func ClientP2PConfig(tx *query.Query, p2pCode string) {
	p2p, _ := tx.GostClientP2P.Preload(tx.GostClientP2P.Node).Where(tx.GostClientP2P.Code.Eq(p2pCode)).First()
	if p2p == nil {
		return
	}
	engine, ok := EngineRegistry.Get(p2p.ClientCode)
	if !ok {
		return
	}
	engine.GetClient().P2PConfig(tx, p2pCode)
}

func ClientRemoveP2PConfig(p2p model.GostClientP2P) {
	engine, ok := EngineRegistry.Get(p2p.ClientCode)
	if !ok {
		return
	}
	engine.GetClient().RemoveP2P(nil, p2p)
}

func ClientPortCheck(tx *query.Query, code string, port string) error {
	engine, ok := EngineRegistry.Get(code)
	if !ok {
		return nil
	}
	return engine.GetClient().PortCheck(tx, "", port)
}
