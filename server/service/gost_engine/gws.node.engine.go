package gost_engine

import (
	"encoding/json"
	"fmt"
	"github.com/go-gost/x/config"
	"github.com/google/uuid"
	"github.com/lxzan/gws"
	"net/http"
	"server/model"
	v1 "server/pkg/p2p_cfg/v1"
	"server/repository"
	"server/repository/query"
	"server/service/common/cache"
	"server/service/common/warn_msg"
	"time"
)

type GwsNodeEngine struct {
	ip         string
	code       string
	conn       *gws.Conn
	registered bool
}

func (e *GwsNodeEngine) checkRegister() {
	time.Sleep(time.Second * 10)
	if !e.registered {
		e.Stop("节点连接注册超时")
	}
	//e.log.Warn("节点连接注册超时，关闭连接", zap.String("type", "node"), zap.String("code", event.code))
}

func (e *GwsNodeEngine) OnOpen(socket *gws.Conn) {
	cache.SetNodeOnline(e.code, true, time.Second*30)
	go func() {
		for {
			time.Sleep(time.Second * 10)
			if err := socket.WritePing(nil); err != nil {
				return
			}
		}
	}()
	go e.checkRegister()
}

func (e *GwsNodeEngine) OnClose(socket *gws.Conn, err error) {
	cache.SetNodeOnline(e.code, false, time.Second*30)
}

func (e *GwsNodeEngine) OnPing(socket *gws.Conn, payload []byte) {
	cache.SetNodeOnline(e.code, true, time.Second*30)
	_ = socket.SetDeadline(time.Now().Add(time.Second * 30))
}

func (e *GwsNodeEngine) OnPong(socket *gws.Conn, payload []byte) {

}

type NodePortCheckResult struct {
	Code string `json:"code"` // 节点编号
	Use  bool   `json:"use"`  // 1=被占用
	Port string `json:"port"` // 端口
}

func (e *GwsNodeEngine) OnMessage(socket *gws.Conn, message *gws.Message) {
	var msg MessageData
	_ = json.Unmarshal(message.Bytes(), &msg)
	switch msg.OperationType {
	case "register":
		e.registered = true
		var data = make(map[string]any)
		_ = msg.GetContent(&data)
		cache.SetNodeVersion(e.code, data["version"].(string))

		_, customDomain := data["custom_domain"]
		_, enableCache := data["cache"]
		cache.SetNodeCustomDomain(e.code, customDomain)
		cache.SetNodeCache(e.code, enableCache)
		db, _, _ := repository.Get("")
		e.Config(db)
	case "port_check":
		var data NodePortCheckResult
		_ = msg.GetContent(&data)
		if data.Code == "" {
			return
		}
		cache.SetNodePortUse(e.code, data.Port, data.Use, time.Minute)
	}
}

func (e *GwsNodeEngine) Stop(msg string) {
	e.WriteMessage(NewMessage(uuid.NewString(), "stop", map[string]string{
		"reason": msg,
	}))
}

func (e *GwsNodeEngine) PortCheck(tx *query.Query, ip, port string) (async bool, allowUse bool) {
	node, err := tx.GostNode.Where(tx.GostNode.Code.Eq(e.code)).First()
	if err != nil {
		return
	}
	var baseConfig model.SystemConfigBase
	cache.GetSystemConfigBase(&baseConfig)
	data := node.GenerateNodePortCheck(baseConfig.BaseUrl, port)
	e.WriteMessage(NewMessage(uuid.NewString(), "port_check", data))
	return false, false
}

type NodeConfigData struct {
	SvcList    []config.ServiceConfig
	Auther     config.AutherConfig
	Ingress    config.IngressConfig
	Limiter    config.LimiterConfig
	Obs        config.ObserverConfig
	P2PCfgCode string
	P2PCfg     v1.ServerConfig
}

func (e *GwsNodeEngine) Config(tx *query.Query) {
	node, err := tx.GostNode.Where(tx.GostNode.Code.Eq(e.code)).First()
	if err != nil {
		return
	}

	var baseConfig model.SystemConfigBase
	cache.GetSystemConfigBase(&baseConfig)

	var data NodeConfigData
	auther := node.GenerateAuther(baseConfig.BaseUrl)
	hosts, _ := tx.GostClientHost.Where(tx.GostClientHost.NodeCode.Eq(node.Code)).Find()
	tunnels, _ := tx.GostClientTunnel.Where(tx.GostClientTunnel.NodeCode.Eq(node.Code)).Find()
	ingress := node.GenerateIngress(hosts, tunnels, cache.GetNodeCustomDomain(e.code))
	limiter := node.GenerateLimiter(baseConfig.BaseUrl)
	p2pCfg := node.GenerateP2PServiceConfig(baseConfig.BaseUrl)
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
	if node.P2P == 1 {
		data.P2PCfgCode = node.Code
		data.P2PCfg = p2pCfg
	}
	e.WriteMessage(NewMessage(uuid.NewString(), "config", data))
}

func (e *GwsNodeEngine) Ingress(tx *query.Query) {
	node, _ := tx.GostNode.Where(tx.GostNode.Code.Eq(e.code)).First()
	if node == nil {
		return
	}

	hosts, _ := tx.GostClientHost.Where(tx.GostClientHost.NodeCode.Eq(e.code)).Find()
	var newHosts []*model.GostClientHost
	for _, host := range hosts {
		if warn_msg.GetHostWarnMsg(*host) != "" {
			continue
		}
		newHosts = append(newHosts, host)
	}

	tunnels, _ := tx.GostClientTunnel.Where(tx.GostClientTunnel.NodeCode.Eq(e.code)).Find()
	var newTunnels []*model.GostClientTunnel
	for _, tunnel := range tunnels {
		if warn_msg.GetTunnelWarnMsg(*tunnel) != "" {
			continue
		}
		newTunnels = append(newTunnels, tunnel)
	}
	var data NodeConfigData
	data.Ingress = node.GenerateIngress(newHosts, newTunnels, cache.GetNodeCustomDomain(e.code))
	e.WriteMessage(NewMessage(uuid.NewString(), "config", data))
}

type HttpsDomainData struct {
	Domain     string
	Target     string
	Cert       string
	Key        string
	ForceHttps int
}

func (e *GwsNodeEngine) CustomDomain(tx *query.Query, domain, cert, key string, forceHttps int) {
	node, err := tx.GostNode.Where(tx.GostNode.Code.Eq(e.code)).First()
	if err != nil {
		return
	}
	e.WriteMessage(NewMessage(uuid.NewString(), "https_domain", HttpsDomainData{
		Domain:     domain,
		Target:     fmt.Sprintf("http://127.0.0.1:%s", node.TunnelInPort),
		Cert:       cert,
		Key:        key,
		ForceHttps: forceHttps,
	}))
}

func NewGwsNodeEngine(code, ip string, w http.ResponseWriter, r *http.Request) (*GwsNodeEngine, error) {
	engine := GwsNodeEngine{
		ip:   ip,
		code: code,
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

func (e *GwsNodeEngine) WriteMessage(req MessageData) {
	marshal, _ := json.Marshal(req)
	_ = e.conn.WriteString(string(marshal))
}
