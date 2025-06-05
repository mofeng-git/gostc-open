package gost_engine

import (
	"server/model"
	"server/repository/query"
	"sync"
)

type IClient interface {
	Stop(msg string)                                                                  // 停止运行
	PortCheck(tx *query.Query, ip, port string) (async bool, allowUse bool)           // 检测端口可用性
	HostConfig(tx *query.Query, hostCode string)                                      // 域名解析配置
	RemoveHost(tx *query.Query, host model.GostClientHost, node model.GostNode)       // 删除域名解析
	ForwardConfig(tx *query.Query, forwardCode string)                                // 端口转发配置
	RemoveForward(tx *query.Query, forward model.GostClientForward)                   // 删除端口转发
	TunnelConfig(tx *query.Query, tunnelCode string)                                  // 私有隧道配置
	RemoveTunnel(tx *query.Query, tunnel model.GostClientTunnel, node model.GostNode) // 删除私有隧道
	P2PConfig(tx *query.Query, p2pCode string)                                        // 私有隧道配置
	RemoveP2P(tx *query.Query, p2p model.GostClientP2P)                               // 删除P2P隧道
	ProxyConfig(tx *query.Query, proxyCode string)                                    // 代理隧道配置
	RemoveProxy(tx *query.Query, proxy model.GostClientProxy)                         // 删除代理隧道
}

type INode interface {
	Stop(msg string)                                                        // 停止运行
	PortCheck(tx *query.Query, ip, port string) (async bool, allowUse bool) // 检测端口可用性
	Config(tx *query.Query)                                                 // 服务端配置
	Ingress(tx *query.Query)                                                // 域名解析映射配置
	CustomDomain(tx *query.Query, domain, cert, key string, forceHttps int) // 自定义域名
}

type Engine struct {
	code       string
	clientImpl IClient
	nodeImpl   INode
}

func (e *Engine) GetClient() IClient {
	return e.clientImpl
}
func (e *Engine) GetNode() INode {
	return e.nodeImpl
}

func NewClientEngine(code string, impl IClient) *Engine {
	return &Engine{
		code:       code,
		clientImpl: impl,
		nodeImpl:   &DefaultNodeImpl{},
	}
}

func NewNodeEngine(code string, impl INode) *Engine {
	return &Engine{
		code:       code,
		clientImpl: &DefaultClientImpl{},
		nodeImpl:   impl,
	}
}

var EngineRegistry = NewEngineRegister()

func NewEngineRegister() *EngineRegister {
	return &EngineRegister{
		data: make(map[string]*Engine),
		m:    &sync.RWMutex{},
	}
}

type EngineRegister struct {
	data map[string]*Engine
	m    *sync.RWMutex
}

func (reg EngineRegister) Set(engine *Engine) {
	reg.m.Lock()
	defer reg.m.Unlock()
	if engine.code == "" {
		return
	}
	reg.data[engine.code] = engine
}

func (reg EngineRegister) Get(code string) (*Engine, bool) {
	reg.m.RLock()
	defer reg.m.RUnlock()
	engine, ok := reg.data[code]
	return engine, ok
}

type DefaultClientImpl struct {
}

func (d DefaultClientImpl) Stop(msg string) {
}

func (d DefaultClientImpl) PortCheck(tx *query.Query, ip, port string) (async bool, allowUse bool) {
	return false, false
}

func (d DefaultClientImpl) HostConfig(tx *query.Query, hostCode string) {
}

func (d DefaultClientImpl) RemoveHost(tx *query.Query, host model.GostClientHost, node model.GostNode) {
}

func (d DefaultClientImpl) ForwardConfig(tx *query.Query, forwardCode string) {
}

func (d DefaultClientImpl) RemoveForward(tx *query.Query, forward model.GostClientForward) {
}

func (d DefaultClientImpl) TunnelConfig(tx *query.Query, tunnelCode string) {
}

func (d DefaultClientImpl) RemoveTunnel(tx *query.Query, tunnel model.GostClientTunnel, node model.GostNode) {
}

func (d DefaultClientImpl) P2PConfig(tx *query.Query, p2pCode string) {
}

func (d DefaultClientImpl) RemoveP2P(tx *query.Query, p2p model.GostClientP2P) {
}

func (d DefaultClientImpl) ProxyConfig(tx *query.Query, proxyCode string) {
}

func (d DefaultClientImpl) RemoveProxy(tx *query.Query, proxy model.GostClientProxy) {
}

func (d DefaultClientImpl) TunConfig(tx *query.Query, tunCode string) {
}

func (d DefaultClientImpl) RemoveTun(tx *query.Query, tunCode string) {
}

type DefaultNodeImpl struct {
}

func (d DefaultNodeImpl) Stop(msg string) {
}

func (d DefaultNodeImpl) PortCheck(tx *query.Query, ip, port string) (async bool, allowUse bool) {
	return false, false
}

func (d DefaultNodeImpl) Config(tx *query.Query) {
}

func (d DefaultNodeImpl) Ingress(tx *query.Query) {
}

func (d DefaultNodeImpl) CustomDomain(tx *query.Query, domain, cert, key string, forceHttps int) {
}

func (d DefaultNodeImpl) DomainCache(tx *query.Query, hostCode string) {
}

func (d DefaultNodeImpl) DomainCacheClean(tx *query.Query, hostCode string) {
}
