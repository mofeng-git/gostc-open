package engine

import (
	"server/model"
	"server/repository/cache"
	"server/repository/query"
)

func ClientStop(code string, msg string) {
	engine, ok := EngineRegistry.Get(code)
	if !ok {
		return
	}
	if !engine.GetClient().IsRunning() {
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
	if !engine.GetClient().IsRunning() {
		return
	}
	engine.GetClient().HostConfig(tx, hostCode)
}

func ClientRemoveHostConfig(tx *query.Query, host model.GostClientHost, node model.GostNode) {
	engine, ok := EngineRegistry.Get(host.ClientCode)
	if !ok {
		return
	}
	if !engine.GetClient().IsRunning() {
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
	if !engine.GetClient().IsRunning() {
		return
	}
	engine.GetClient().ForwardConfig(tx, forwardCode)
}

func ClientRemoveForwardConfig(forward model.GostClientForward) {
	engine, ok := EngineRegistry.Get(forward.ClientCode)
	if !ok {
		return
	}
	if !engine.GetClient().IsRunning() {
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
	if !engine.GetClient().IsRunning() {
		return
	}
	engine.GetClient().TunnelConfig(tx, tunnelCode)
}

func ClientRemoveTunnelConfig(tx *query.Query, tunnel model.GostClientTunnel, node model.GostNode) {
	engine, ok := EngineRegistry.Get(tunnel.ClientCode)
	if !ok {
		return
	}
	if !engine.GetClient().IsRunning() {
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
	if !engine.GetClient().IsRunning() {
		return
	}
	engine.GetClient().ProxyConfig(tx, proxyCode)
}

func ClientRemoveProxyConfig(proxy model.GostClientProxy) {
	engine, ok := EngineRegistry.Get(proxy.ClientCode)
	if !ok {
		return
	}
	if !engine.GetClient().IsRunning() {
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
	if !engine.GetClient().IsRunning() {
		return
	}
	engine.GetClient().P2PConfig(tx, p2pCode)
}

func ClientRemoveP2PConfig(p2p model.GostClientP2P) {
	engine, ok := EngineRegistry.Get(p2p.ClientCode)
	if !ok {
		return
	}
	if !engine.GetClient().IsRunning() {
		return
	}
	engine.GetClient().RemoveP2P(nil, p2p)
}

func ClientPortCheck(tx *query.Query, code string, port string) error {
	engine, ok := EngineRegistry.Get(code)
	if !ok {
		return nil
	}
	if !engine.GetClient().IsRunning() {
		return nil
	}
	return engine.GetClient().PortCheck(tx, "", port)
}

func ClientCfgConfig(tx *query.Query, code string) error {
	cfg, err := tx.FrpClientCfg.Where(tx.FrpClientCfg.Code.Eq(code)).First()
	if err != nil || cfg == nil {
		return nil
	}
	engine, ok := EngineRegistry.Get(cfg.ClientCode)
	if !ok {
		return nil
	}
	if !engine.GetClient().IsRunning() {
		return nil
	}
	return engine.GetClient().CustomCfgConfig(tx, code)
}

func ClientCfgRemove(clientCode, code string) {
	engine, ok := EngineRegistry.Get(clientCode)
	if !ok {
		return
	}
	if !engine.GetClient().IsRunning() {
		return
	}
	_ = engine.GetClient().RemoveCustomCfg(nil, code)
}

// 将客户端的所有隧道下发到客户端上
func ClientAllConfigUpdateByClientCode(tx *query.Query, clientCode string) {
	if !cache.GetClientOnline(clientCode) {
		return
	}

	var hostCodes []string
	_ = tx.GostClientHost.Where(
		tx.GostClientHost.ClientCode.Eq(clientCode),
	).Pluck(tx.GostClientHost.Code, &hostCodes)
	for _, code := range hostCodes {
		ClientHostConfig(tx, code)
	}
	var forwardCodes []string
	_ = tx.GostClientForward.Where(
		tx.GostClientForward.ClientCode.Eq(clientCode),
	).Pluck(tx.GostClientForward.Code, &forwardCodes)
	for _, code := range forwardCodes {
		ClientForwardConfig(tx, code)
	}
	var tunnelCodes []string
	_ = tx.GostClientTunnel.Where(
		tx.GostClientTunnel.ClientCode.Eq(clientCode),
	).Pluck(tx.GostClientTunnel.Code, &tunnelCodes)
	for _, code := range tunnelCodes {
		ClientTunnelConfig(tx, code)
	}
	var proxyCodes []string
	_ = tx.GostClientProxy.Where(
		tx.GostClientProxy.ClientCode.Eq(clientCode),
	).Pluck(tx.GostClientProxy.Code, &proxyCodes)
	for _, code := range proxyCodes {
		ClientProxyConfig(tx, code)
	}
	var p2pCodes []string
	_ = tx.GostClientP2P.Where(
		tx.GostClientP2P.ClientCode.Eq(clientCode),
	).Pluck(tx.GostClientP2P.Code, &p2pCodes)
	for _, code := range p2pCodes {
		ClientP2PConfig(tx, code)
	}
	var cfgsCode []string
	_ = tx.FrpClientCfg.Where(
		tx.FrpClientCfg.ClientCode.Eq(clientCode),
	).Pluck(tx.FrpClientCfg.Code, &cfgsCode)
	for _, code := range cfgsCode {
		ClientCfgConfig(tx, code)
	}
}

// 将节点的所有隧道下发到客户端上
func ClientAllConfigUpdateByNodeCode(tx *query.Query, nodeCode string) {
	if !cache.GetNodeOnline(nodeCode) {
		return
	}

	// 查询客户端，只过滤出在线的客户端，减少多余数据处理
	var clientCodes []string
	_ = tx.GostClient.Pluck(tx.GostClient.Code, &clientCodes)
	var onlineClientCodes []string
	for _, code := range clientCodes {
		if cache.GetClientOnline(code) {
			onlineClientCodes = append(onlineClientCodes, code)
		}
	}

	var hostCodes []string
	_ = tx.GostClientHost.Where(
		tx.GostClientHost.NodeCode.Eq(nodeCode),
	).Pluck(tx.GostClientHost.Code, &hostCodes)
	for _, code := range hostCodes {
		ClientHostConfig(tx, code)
	}
	var forwardCodes []string
	_ = tx.GostClientForward.Where(
		tx.GostClientForward.NodeCode.Eq(nodeCode),
	).Pluck(tx.GostClientForward.Code, &forwardCodes)
	for _, code := range forwardCodes {
		ClientForwardConfig(tx, code)
	}
	var tunnelCodes []string
	_ = tx.GostClientTunnel.Where(
		tx.GostClientTunnel.NodeCode.Eq(nodeCode),
	).Pluck(tx.GostClientTunnel.Code, &tunnelCodes)
	for _, code := range tunnelCodes {
		ClientTunnelConfig(tx, code)
	}
	var proxyCodes []string
	_ = tx.GostClientProxy.Where(
		tx.GostClientProxy.NodeCode.Eq(nodeCode),
	).Pluck(tx.GostClientProxy.Code, &proxyCodes)
	for _, code := range proxyCodes {
		ClientProxyConfig(tx, code)
	}
	var p2pCodes []string
	_ = tx.GostClientP2P.Where(
		tx.GostClientP2P.NodeCode.Eq(nodeCode),
	).Pluck(tx.GostClientP2P.Code, &p2pCodes)
	for _, code := range p2pCodes {
		ClientP2PConfig(tx, code)
	}
}
