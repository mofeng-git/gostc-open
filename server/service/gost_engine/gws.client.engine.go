package gost_engine

import (
	"encoding/json"
	"github.com/go-gost/x/config"
	"github.com/google/uuid"
	"github.com/lxzan/gws"
	cache2 "github.com/patrickmn/go-cache"
	"net/http"
	"server/model"
	v1 "server/pkg/p2p_cfg/v1"
	"server/repository"
	"server/repository/query"
	"server/service/common/cache"
	"server/service/common/warn_msg"
	"time"
)

type GwsClientEngine struct {
	code       string
	conn       *gws.Conn
	registered bool
}

func (e *GwsClientEngine) checkRegister() {
	time.Sleep(time.Second * 10)
	if !e.registered {
		e.Stop("客户端连接注册超时")
	}
	//e.log.Warn("客户端连接注册超时，关闭连接", zap.String("type", "client"), zap.String("code", event.code))
}

func (e *GwsClientEngine) register() {
	e.registered = true
	db, _, _ := repository.Get("")
	var hostCodes []string
	_ = db.GostClientHost.Where(db.GostClientHost.ClientCode.Eq(e.code)).Pluck(db.GostClientHost.Code, &hostCodes)
	for _, code := range hostCodes {
		e.HostConfig(db, code)
	}
	var forwardCodes []string
	_ = db.GostClientForward.Where(db.GostClientForward.ClientCode.Eq(e.code)).Pluck(db.GostClientForward.Code, &forwardCodes)
	for _, code := range forwardCodes {
		e.ForwardConfig(db, code)
	}
	var tunnelCodes []string
	_ = db.GostClientTunnel.Where(db.GostClientTunnel.ClientCode.Eq(e.code)).Pluck(db.GostClientTunnel.Code, &tunnelCodes)
	for _, code := range tunnelCodes {
		e.TunnelConfig(db, code)
	}
	var proxyCodes []string
	_ = db.GostClientProxy.Where(db.GostClientProxy.ClientCode.Eq(e.code)).Pluck(db.GostClientProxy.Code, &proxyCodes)
	for _, code := range proxyCodes {
		e.ProxyConfig(db, code)
	}
	var p2pCodes []string
	_ = db.GostClientP2P.Where(db.GostClientP2P.ClientCode.Eq(e.code)).Pluck(db.GostClientP2P.Code, &p2pCodes)
	for _, code := range p2pCodes {
		e.P2PConfig(db, code)
	}
}

func (e *GwsClientEngine) OnOpen(socket *gws.Conn) {
	go func() {
		for {
			time.Sleep(time.Second * 10)
			if err := socket.WritePing(nil); err != nil {
				return
			}
		}
	}()
	cache.SetClientOnline(e.code, true, time.Second*30)
	cache.SetClientLastTime(e.code)
	go e.checkRegister()
}

func (e *GwsClientEngine) OnClose(socket *gws.Conn, err error) {
	cache.SetClientOnline(e.code, false, time.Second*30)
}

func (e *GwsClientEngine) OnPing(socket *gws.Conn, payload []byte) {
	cache.SetClientOnline(e.code, true, time.Second*30)
	_ = socket.SetDeadline(time.Now().Add(time.Second * 30))
}

func (e *GwsClientEngine) OnPong(socket *gws.Conn, payload []byte) {

}

type ClientPortCheckResult struct {
	Code string `json:"code"` // 节点编号
	Use  bool   `json:"use"`  // 1=被占用
	Port string `json:"port"` // 端口
}

func (e *GwsClientEngine) OnMessage(socket *gws.Conn, message *gws.Message) {
	var msg MessageData
	_ = json.Unmarshal(message.Bytes(), &msg)
	switch msg.OperationType {
	case "register":
		var data = make(map[string]any)
		_ = msg.GetContent(&data)
		cache.SetClientVersion(e.code, data["version"].(string))
		e.register()
	case "port_check":
		var data ClientPortCheckResult
		_ = msg.GetContent(&data)
		if data.Code == "" {
			return
		}
		cache.SetClientPortUse(e.code, data.Port, data.Use, cache2.NoExpiration)
	}
}

func (e *GwsClientEngine) PortCheck(tx *query.Query, ip, port string) (async bool, allowUse bool) {
	client, err := tx.GostClient.Where(tx.GostClient.Code.Eq(e.code)).First()
	if err != nil {
		return
	}

	var baseConfig model.SystemConfigBase
	cache.GetSystemConfigBase(&baseConfig)
	data := client.GenerateClientPortCheck(baseConfig.BaseUrl, port)
	e.WriteMessage(NewMessage(uuid.NewString(), "port_check", data))
	return false, false
}

func (e *GwsClientEngine) Stop(msg string) {
	e.WriteMessage(NewMessage(uuid.NewString(), "stop", map[string]string{
		"reason": msg,
	}))
}

type ClientHostConfigData struct {
	Svc           config.ServiceConfig
	Chain         config.ChainConfig
	AdmissionList []config.AdmissionConfig
}

func (e *GwsClientEngine) HostConfig(tx *query.Query, hostCode string) {
	host, _ := tx.GostClientHost.Preload(tx.GostClientHost.Node).Where(tx.GostClientHost.Code.Eq(hostCode)).First()
	if host == nil {
		return
	}
	if warn_msg.GetHostWarnMsg(*host) != "" {
		e.RemoveHost(tx, *host, host.Node)
		return
	}
	auth, _ := tx.GostAuth.Where(tx.GostAuth.TunnelCode.Eq(hostCode)).First()
	if auth == nil {
		return
	}
	var data ClientHostConfigData
	chain := host.GenerateChainConfig(*auth)
	admissionWhite, admissionBlack := host.GenerateWhiteAdmission(), host.GenerateBlackAdmission()
	svcCfg, ok := host.GenerateSvcConfig(chain.Name, admissionWhite.Name, admissionBlack.Name)
	if !ok {
		return
	}
	data.Chain = chain
	data.Svc = svcCfg
	data.AdmissionList = []config.AdmissionConfig{admissionWhite, admissionBlack}
	cache.SetIngress(host.DomainPrefix+"."+host.Node.Domain, host.Code)
	cache.SetIngress(host.DomainPrefix+"."+host.Node.Domain+":"+host.Node.TunnelInPort, host.Code)
	cache.SetIngress(host.DomainPrefix+"."+host.Node.Domain+":80", host.Code)
	cache.SetIngress(host.DomainPrefix+"."+host.Node.Domain+":443", host.Code)
	e.WriteMessage(NewMessage(uuid.NewString(), "host_config", data))
	NodeIngress(tx, host.NodeCode)
}

func (e *GwsClientEngine) RemoveHost(tx *query.Query, host model.GostClientHost, node model.GostNode) {
	cache.DelIngress(host.DomainPrefix + "." + node.Domain)
	cache.DelIngress(host.DomainPrefix + "." + node.Domain + ":" + node.TunnelInPort)
	cache.DelIngress(host.DomainPrefix + "." + node.Domain + ":80")
	cache.DelIngress(host.DomainPrefix + "." + node.Domain + ":443")
	e.WriteMessage(NewMessage(uuid.NewString(), "remove_config", []string{host.Code}))
	NodeIngress(tx, node.Code)
}

type ClientForwardConfigData struct {
	SvcList       []config.ServiceConfig
	Chain         config.ChainConfig
	Limiter       config.LimiterConfig
	CLimiter      config.LimiterConfig
	RLimiter      config.LimiterConfig
	Obs           config.ObserverConfig
	AdmissionList []config.AdmissionConfig
}

func (e *GwsClientEngine) ForwardConfig(tx *query.Query, forwardCode string) {
	forward, _ := tx.GostClientForward.Preload(tx.GostClientForward.Node).Where(tx.GostClientForward.Code.Eq(forwardCode)).First()
	if forward == nil {
		return
	}
	if warn_msg.GetForwardWarnMsg(*forward) != "" {
		e.RemoveForward(tx, *forward)
		return
	}

	auth, _ := tx.GostAuth.Where(tx.GostAuth.TunnelCode.Eq(forwardCode)).First()
	if auth == nil {
		return
	}

	var baseConfig model.SystemConfigBase
	cache.GetSystemConfigBase(&baseConfig)
	nodeVersion := cache.GetNodeVersion(forward.NodeCode)

	var data ClientForwardConfigData
	chain := forward.GenerateChainConfig(*auth)
	limiter := forward.GenerateLimiter()
	rLimiter := forward.GenerateRLimiter()
	cLimiter := forward.GenerateCLimiter()
	obs := forward.GenerateObs(baseConfig.BaseUrl, nodeVersion)
	admissionWhite, admissionBlack := forward.GenerateWhiteAdmission(), forward.GenerateBlackAdmission()

	tcpSvcCfg, ok := forward.GenerateTcpSvcConfig(chain.Name, limiter.Name, cLimiter.Name, rLimiter.Name, obs.Name, admissionWhite.Name, admissionBlack.Name)
	if ok {
		data.SvcList = append(data.SvcList, tcpSvcCfg)
	}
	udpSvcCfg, ok := forward.GenerateUdpSvcConfig(chain.Name, limiter.Name, cLimiter.Name, rLimiter.Name, obs.Name, admissionWhite.Name, admissionBlack.Name)
	if ok {
		data.SvcList = append(data.SvcList, udpSvcCfg)
	}
	data.Chain = chain
	data.Limiter = limiter
	data.CLimiter = cLimiter
	data.RLimiter = rLimiter
	data.Obs = obs
	data.AdmissionList = []config.AdmissionConfig{admissionWhite, admissionBlack}
	e.WriteMessage(NewMessage(uuid.NewString(), "forward_config", data))
}

func (e *GwsClientEngine) RemoveForward(tx *query.Query, forward model.GostClientForward) {
	e.WriteMessage(NewMessage(uuid.NewString(), "remove_config", []string{
		"tcp_" + forward.Code,
		"udp_" + forward.Code,
	}))
}

type ClientTunnelConfigData struct {
	SvcList []config.ServiceConfig
	Chain   config.ChainConfig
}

func (e *GwsClientEngine) TunnelConfig(tx *query.Query, tunnelCode string) {
	tunnel, _ := tx.GostClientTunnel.Preload(tx.GostClientTunnel.Node).Where(tx.GostClientTunnel.Code.Eq(tunnelCode)).First()
	if tunnel == nil {
		return
	}
	if warn_msg.GetTunnelWarnMsg(*tunnel) != "" {
		e.RemoveTunnel(tx, *tunnel, tunnel.Node)
		return
	}

	auth, _ := tx.GostAuth.Where(tx.GostAuth.TunnelCode.Eq(tunnelCode)).First()
	if auth == nil {
		return
	}

	var data ClientTunnelConfigData
	chain := tunnel.GenerateChainConfig(*auth)
	tcpSvcCfg, ok := tunnel.GenerateTcpSvcConfig(chain.Name)
	if ok {
		data.SvcList = append(data.SvcList, tcpSvcCfg)
	}
	udpSvcCfg, ok := tunnel.GenerateUdpSvcConfig(chain.Name)
	if ok {
		data.SvcList = append(data.SvcList, udpSvcCfg)
	}
	data.Chain = chain
	cache.SetIngress(tunnel.Code, tunnel.Code)
	e.WriteMessage(NewMessage(uuid.NewString(), "tunnel_config", data))
	NodeIngress(tx, tunnel.NodeCode)
}

func (e *GwsClientEngine) RemoveTunnel(tx *query.Query, tunnel model.GostClientTunnel, node model.GostNode) {
	cache.DelIngress(tunnel.Code)
	e.WriteMessage(NewMessage(uuid.NewString(), "remove_config", []string{
		"tcp_" + tunnel.Code,
		"udp_" + tunnel.Code,
	}))
	NodeIngress(tx, node.Code)
}

type ClientProxyConfigData struct {
	Svc      config.ServiceConfig
	Chain    config.ChainConfig
	Limiter  config.LimiterConfig
	CLimiter config.LimiterConfig
	RLimiter config.LimiterConfig
	Obs      config.ObserverConfig
}

func (e *GwsClientEngine) P2PConfig(tx *query.Query, p2pCode string) {
	p2p, _ := tx.GostClientP2P.Preload(tx.GostClientP2P.Node).Where(tx.GostClientP2P.Code.Eq(p2pCode)).First()
	if p2p == nil {
		return
	}
	if warn_msg.GetP2PWarnMsg(*p2p) != "" {
		e.RemoveP2P(tx, *p2p)
		return
	}

	var data ClientP2PConfigData
	if p2p.Node.P2P == 1 {
		data.Code = p2pCode
		data.Common, _ = p2p.GenerateCommonCfg()
		data.STCPCfg, data.XTCPCfg = p2p.GenerateProxyCfgs()
	}
	e.WriteMessage(NewMessage(uuid.NewString(), "p2p_config", data))
}

func (e *GwsClientEngine) RemoveP2P(tx *query.Query, p2p model.GostClientP2P) {
	e.WriteMessage(NewMessage(uuid.NewString(), "remove_config", []string{
		p2p.Code,
	}))
}

type ClientP2PConfigData struct {
	Code    string
	Common  v1.ClientCommonConfig
	STCPCfg v1.STCPProxyConfig
	XTCPCfg v1.XTCPProxyConfig
}

func (e *GwsClientEngine) ProxyConfig(tx *query.Query, proxyCode string) {
	proxy, _ := tx.GostClientProxy.Preload(tx.GostClientProxy.Node).Where(tx.GostClientProxy.Code.Eq(proxyCode)).First()
	if proxy == nil {
		return
	}
	if warn_msg.GetProxyWarnMsg(*proxy) != "" {
		e.RemoveProxy(tx, *proxy)
		return
	}

	auth, _ := tx.GostAuth.Where(tx.GostAuth.TunnelCode.Eq(proxyCode)).First()
	if auth == nil {
		return
	}

	var baseConfig model.SystemConfigBase
	cache.GetSystemConfigBase(&baseConfig)
	nodeVersion := cache.GetNodeVersion(proxy.NodeCode)

	var data ClientProxyConfigData
	chain := proxy.GenerateChainConfig(*auth)
	limiter := proxy.GenerateLimiter()
	rLimiter := proxy.GenerateRLimiter()
	cLimiter := proxy.GenerateCLimiter()
	obs := proxy.GenerateObs(baseConfig.BaseUrl, nodeVersion)

	svcCfg, ok := proxy.GenerateSvcConfig(chain.Name, limiter.Name, cLimiter.Name, rLimiter.Name, obs.Name)
	if ok {
		data.Svc = svcCfg
	}
	data.Chain = chain
	data.Limiter = limiter
	data.CLimiter = cLimiter
	data.RLimiter = rLimiter
	data.Obs = obs
	e.WriteMessage(NewMessage(uuid.NewString(), "proxy_config", data))
}

func (e *GwsClientEngine) RemoveProxy(tx *query.Query, proxy model.GostClientProxy) {
	e.WriteMessage(NewMessage(uuid.NewString(), "remove_config", []string{
		"proxy_" + proxy.Code,
	}))
}

func NewGwsClientEngine(code string, w http.ResponseWriter, r *http.Request) (*GwsClientEngine, error) {
	engine := GwsClientEngine{
		code: code,
		conn: nil,
	}
	upgrader := gws.NewUpgrader(&engine, &gws.ServerOption{
		ParallelEnabled:   true,                                 // 开启并行消息处理
		Recovery:          gws.Recovery,                         // 开启异常恢复
		PermessageDeflate: gws.PermessageDeflate{Enabled: true}, // 开启压缩
	})
	upgrade, err := upgrader.Upgrade(w, r)
	if err != nil {
		return nil, err
	}
	engine.conn = upgrade
	go engine.conn.ReadLoop()
	return &engine, nil
}

func (e *GwsClientEngine) WriteMessage(req MessageData) {
	marshal, _ := json.Marshal(req)
	_ = e.conn.WriteString(string(marshal))
}
