package engine

import (
	"github.com/lesismal/arpc/codec"
	"gopkg.in/yaml.v3"
	"server/model"
	"server/repository/query"
	"sync"
)

func init() {
	codec.SetCodec(&YAMLCodec{})
}

type YAMLCodec struct {
}

func (y *YAMLCodec) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (y *YAMLCodec) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

type IClient interface {
	IsRunning() bool                                                                        // 运行状态
	Stop(msg string)                                                                        // 停止运行
	PortCheck(tx *query.Query, ip, port string) error                                       // 检测端口可用性
	HostConfig(tx *query.Query, hostCode string) error                                      // 域名解析配置
	RemoveHost(tx *query.Query, host model.GostClientHost, node model.GostNode) error       // 删除域名解析
	ForwardConfig(tx *query.Query, forwardCode string) error                                // 端口转发配置
	RemoveForward(tx *query.Query, forward model.GostClientForward) error                   // 删除端口转发
	TunnelConfig(tx *query.Query, tunnelCode string) error                                  // 私有隧道配置
	RemoveTunnel(tx *query.Query, tunnel model.GostClientTunnel, node model.GostNode) error // 删除私有隧道
	P2PConfig(tx *query.Query, p2pCode string) error                                        // 私有隧道配置
	RemoveP2P(tx *query.Query, p2p model.GostClientP2P) error                               // 删除P2P隧道
	ProxyConfig(tx *query.Query, proxyCode string) error                                    // 代理隧道配置
	RemoveProxy(tx *query.Query, proxy model.GostClientProxy) error                         // 删除代理隧道
	CustomCfgConfig(tx *query.Query, cfgCode string) error                                  // 自定义配置隧道
	RemoveCustomCfg(tx *query.Query, cfgCode string) error                                  // 删除自定义配置
}

type INode interface {
	IsRunning() bool                                                              // 运行状态
	Stop(msg string)                                                              // 停止运行
	PortCheck(tx *query.Query, ip, port string) error                             // 检测端口可用性
	Config(tx *query.Query) error                                                 // 服务端配置
	Ingress(tx *query.Query) error                                                // 域名解析映射配置
	CustomDomain(tx *query.Query, domain, cert, key string, forceHttps int) error // 自定义域名
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

func (d DefaultClientImpl) IsRunning() bool {
	return false
}

func (d DefaultClientImpl) CustomCfgConfig(tx *query.Query, cfgCode string) error {
	return nil
}

func (d DefaultClientImpl) RemoveCustomCfg(tx *query.Query, cfgCode string) error {
	return nil
}

func (d DefaultClientImpl) Stop(msg string) {
}

func (d DefaultClientImpl) PortCheck(tx *query.Query, ip, port string) error {
	return nil
}

func (d DefaultClientImpl) HostConfig(tx *query.Query, hostCode string) error {
	return nil
}

func (d DefaultClientImpl) RemoveHost(tx *query.Query, host model.GostClientHost, node model.GostNode) error {
	return nil
}

func (d DefaultClientImpl) ForwardConfig(tx *query.Query, forwardCode string) error {
	return nil
}

func (d DefaultClientImpl) RemoveForward(tx *query.Query, forward model.GostClientForward) error {
	return nil
}

func (d DefaultClientImpl) TunnelConfig(tx *query.Query, tunnelCode string) error {
	return nil
}

func (d DefaultClientImpl) RemoveTunnel(tx *query.Query, tunnel model.GostClientTunnel, node model.GostNode) error {
	return nil
}

func (d DefaultClientImpl) P2PConfig(tx *query.Query, p2pCode string) error {
	return nil
}

func (d DefaultClientImpl) RemoveP2P(tx *query.Query, p2p model.GostClientP2P) error {
	return nil
}

func (d DefaultClientImpl) ProxyConfig(tx *query.Query, proxyCode string) error {
	return nil
}

func (d DefaultClientImpl) RemoveProxy(tx *query.Query, proxy model.GostClientProxy) error {
	return nil
}

type DefaultNodeImpl struct {
}

func (d DefaultNodeImpl) IsRunning() bool {
	return false
}

func (d DefaultNodeImpl) Stop(msg string) {
}

func (d DefaultNodeImpl) PortCheck(tx *query.Query, ip, port string) error {
	return nil
}

func (d DefaultNodeImpl) Config(tx *query.Query) error {
	return nil
}

func (d DefaultNodeImpl) Ingress(tx *query.Query) error {
	return nil
}

func (d DefaultNodeImpl) CustomDomain(tx *query.Query, domain, cert, key string, forceHttps int) error {
	return nil
}
