package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-gost/core/logger"
	"github.com/go-gost/x/config"
	"github.com/go-gost/x/config/parsing/admission"
	"github.com/go-gost/x/config/parsing/auth"
	"github.com/go-gost/x/config/parsing/chain"
	"github.com/go-gost/x/config/parsing/ingress"
	"github.com/go-gost/x/config/parsing/limiter"
	"github.com/go-gost/x/config/parsing/observer"
	"github.com/go-gost/x/config/parsing/service"
	"github.com/go-gost/x/registry"
	"github.com/google/uuid"
	"github.com/lxzan/gws"
	"gostc-sub/pkg/utils"
	"net/http"
	"os"
	"strconv"
	"time"
)

var SvcMap = make(map[string]bool)

type Event struct {
	server bool
	key    string
}

func NewEvent(key string, server bool) *Event {
	return &Event{
		server: server,
		key:    key,
	}
}

func (e *Event) OnOpen(socket *gws.Conn) {
	fmt.Println("connect success")
	go func() {
		for {
			time.Sleep(time.Second * 10)
			if err := socket.WritePing(nil); err != nil {
				fmt.Println("send ping msg fail", err)
				_ = socket.WriteClose(1000, nil)
				return
			}
		}
	}()
	version := CLIENT_VERSION
	if e.server {
		version = SERVER_VERSION
	}
	e.WriteAny(socket, NewMessage(uuid.NewString(), "register", map[string]any{
		"version": version,
	}))
}

func (e *Event) OnClose(socket *gws.Conn, err error) {
	if !errors.Is(err, gws.ErrConnClosed) {
		fmt.Println("conn close", err)
	}
}

func (e *Event) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.SetDeadline(time.Now().Add(time.Second * 30))
}

func (e *Event) OnPong(socket *gws.Conn, payload []byte) {
}

type ClientHostConfigData struct {
	Svc           config.ServiceConfig
	Chain         config.ChainConfig
	AdmissionList []config.AdmissionConfig
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

type ClientTunnelConfigData struct {
	SvcList []config.ServiceConfig
	Chain   config.ChainConfig
}

type ClientProxyConfigData struct {
	Svc      config.ServiceConfig
	Chain    config.ChainConfig
	Limiter  config.LimiterConfig
	CLimiter config.LimiterConfig
	RLimiter config.LimiterConfig
	Obs      config.ObserverConfig
}

type ConfigData struct {
	SvcList []config.ServiceConfig
	Auther  config.AutherConfig
	Ingress config.IngressConfig
	Limiter config.LimiterConfig
	Obs     config.ObserverConfig
}

func (e *Event) OnMessage(socket *gws.Conn, message *gws.Message) {
	var msg MessageData
	_ = json.Unmarshal(message.Bytes(), &msg)
	switch msg.OperationType {
	case "stop":
		var data = make(map[string]any)
		_ = msg.GetContent(&data)
		fmt.Println(data["reason"])
		e.WriteAny(socket, NewMessage(msg.OperationId, msg.OperationType, nil))
		time.Sleep(time.Second * 2)
		os.Exit(0)
	case "forward_config":
		var data ClientForwardConfigData
		_ = msg.GetContent(&data)
		parseChain, err := chain.ParseChain(&data.Chain, logger.Default())
		if err == nil {
			registry.ChainRegistry().Unregister(data.Chain.Name)
			_ = registry.ChainRegistry().Register(data.Chain.Name, parseChain)
		}
		if data.Limiter.Name != "" {
			trafficLimiter := limiter.ParseTrafficLimiter(&data.Limiter)
			registry.TrafficLimiterRegistry().Unregister(data.Limiter.Name)
			_ = registry.TrafficLimiterRegistry().Register(data.Limiter.Name, trafficLimiter)
		}
		if data.CLimiter.Name != "" {
			connLimiter := limiter.ParseConnLimiter(&data.CLimiter)
			registry.ConnLimiterRegistry().Unregister(data.CLimiter.Name)
			_ = registry.ConnLimiterRegistry().Register(data.CLimiter.Name, connLimiter)
		}
		if data.RLimiter.Name != "" {
			rateLimiter := limiter.ParseRateLimiter(&data.RLimiter)
			registry.RateLimiterRegistry().Unregister(data.RLimiter.Name)
			_ = registry.RateLimiterRegistry().Register(data.RLimiter.Name, rateLimiter)
		}
		if data.Obs.Name != "" {
			parseObserver := observer.ParseObserver(&data.Obs)
			registry.ObserverRegistry().Unregister(data.Obs.Name)
			_ = registry.ObserverRegistry().Register(data.Obs.Name, parseObserver)
		}

		for _, item := range data.AdmissionList {
			if item.Name != "" {
				parseAdmission := admission.ParseAdmission(&item)
				registry.AdmissionRegistry().Unregister(item.Name)
				_ = registry.AdmissionRegistry().Register(item.Name, parseAdmission)
			}
		}

		for _, svcCfg := range data.SvcList {
			if oldSvc := registry.ServiceRegistry().Get(svcCfg.Name); oldSvc != nil {
				registry.ServiceRegistry().Unregister(svcCfg.Name)
				_ = oldSvc.Close()
			}
			svc, err := service.ParseService(&svcCfg)
			if err != nil {
				continue
			}
			go svc.Serve()
			if err = registry.ServiceRegistry().Register(svcCfg.Name, svc); err != nil {
				_ = svc.Close()
				continue
			}
			SvcMap[svcCfg.Name] = true
		}
		e.WriteAny(socket, NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "host_config":
		var data ClientHostConfigData
		_ = msg.GetContent(&data)
		parseChain, err := chain.ParseChain(&data.Chain, logger.Default())
		if err == nil {
			registry.ChainRegistry().Unregister(data.Chain.Name)
			_ = registry.ChainRegistry().Register(data.Chain.Name, parseChain)
		}

		for _, item := range data.AdmissionList {
			if item.Name != "" {
				parseAdmission := admission.ParseAdmission(&item)
				registry.AdmissionRegistry().Unregister(item.Name)
				_ = registry.AdmissionRegistry().Register(item.Name, parseAdmission)
			}
		}

		if oldSvc := registry.ServiceRegistry().Get(data.Svc.Name); oldSvc != nil {
			registry.ServiceRegistry().Unregister(data.Svc.Name)
			_ = oldSvc.Close()
		}
		svc, err := service.ParseService(&data.Svc)
		if err == nil {
			go svc.Serve()
			if err = registry.ServiceRegistry().Register(data.Svc.Name, svc); err != nil {
				_ = svc.Close()
			}
			SvcMap[data.Svc.Name] = true
		}

		e.WriteAny(socket, NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "tunnel_config":
		var data ClientTunnelConfigData
		_ = msg.GetContent(&data)
		parseChain, err := chain.ParseChain(&data.Chain, logger.Default())
		if err == nil {
			registry.ChainRegistry().Unregister(data.Chain.Name)
			_ = registry.ChainRegistry().Register(data.Chain.Name, parseChain)
		}
		for _, svcCfg := range data.SvcList {
			if oldSvc := registry.ServiceRegistry().Get(svcCfg.Name); oldSvc != nil {
				registry.ServiceRegistry().Unregister(svcCfg.Name)
				_ = oldSvc.Close()
			}
			svc, err := service.ParseService(&svcCfg)
			if err != nil {
				continue
			}
			go svc.Serve()
			if err = registry.ServiceRegistry().Register(svcCfg.Name, svc); err != nil {
				_ = svc.Close()
				continue
			}
			SvcMap[svcCfg.Name] = true
		}
		e.WriteAny(socket, NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "proxy_config":
		var data ClientProxyConfigData
		_ = msg.GetContent(&data)
		parseChain, err := chain.ParseChain(&data.Chain, logger.Default())
		if err == nil {
			registry.ChainRegistry().Unregister(data.Chain.Name)
			_ = registry.ChainRegistry().Register(data.Chain.Name, parseChain)
		}
		if oldSvc := registry.ServiceRegistry().Get(data.Svc.Name); oldSvc != nil {
			registry.ServiceRegistry().Unregister(data.Svc.Name)
			_ = oldSvc.Close()
		}
		svc, err := service.ParseService(&data.Svc)
		if err == nil {
			go svc.Serve()
			if err = registry.ServiceRegistry().Register(data.Svc.Name, svc); err != nil {
				_ = svc.Close()
			}
			SvcMap[data.Svc.Name] = true
		}
		e.WriteAny(socket, NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "remove_config":
		var names []string
		_ = msg.GetContent(&names)
		for _, name := range names {
			if svc := registry.ServiceRegistry().Get(name); svc != nil {
				_ = svc.Close()
				registry.ServiceRegistry().Unregister(name)
				SvcMap[name] = false
			}
		}
		e.WriteAny(socket, NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "config":
		var data ConfigData
		_ = msg.GetContent(&data)
		if data.Limiter.Name != "" {
			trafficLimiter := limiter.ParseTrafficLimiter(&data.Limiter)
			registry.TrafficLimiterRegistry().Unregister(data.Limiter.Name)
			_ = registry.TrafficLimiterRegistry().Register(data.Limiter.Name, trafficLimiter)
		}

		if data.Ingress.Name != "" {
			parseIngress := ingress.ParseIngress(&data.Ingress)
			registry.IngressRegistry().Unregister(data.Ingress.Name)
			_ = registry.IngressRegistry().Register(data.Ingress.Name, parseIngress)
		}

		if data.Auther.Name != "" {
			parseAuther := auth.ParseAuther(&data.Auther)
			registry.AutherRegistry().Unregister(data.Auther.Name)
			_ = registry.AutherRegistry().Register(data.Auther.Name, parseAuther)
		}
		if data.Obs.Name != "" {
			parseObserver := observer.ParseObserver(&data.Obs)
			registry.ObserverRegistry().Unregister(data.Obs.Name)
			_ = registry.ObserverRegistry().Register(data.Obs.Name, parseObserver)
		}

		for _, svcCfg := range data.SvcList {
			if oldSvc := registry.ServiceRegistry().Get(svcCfg.Name); oldSvc != nil {
				registry.ServiceRegistry().Unregister(svcCfg.Name)
				_ = oldSvc.Close()
			}
			svc, err := service.ParseService(&svcCfg)
			if err != nil {
				continue
			}
			go svc.Serve()
			if err = registry.ServiceRegistry().Register(svcCfg.Name, svc); err != nil {
				_ = svc.Close()
				continue
			}
			SvcMap[svcCfg.Name] = true
		}
		e.WriteAny(socket, NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	case "port_check":
		var data = make(map[string]string)
		_ = msg.GetContent(&data)
		url := data["callback"]
		port, _ := strconv.Atoi(data["port"])
		if port != 0 {
			marshal, _ := json.Marshal(PortCheckResp{
				Code: data["code"],
				Use:  utils.IsUse(port),
				Port: data["port"],
			})
			req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(marshal))
			response, err := http.DefaultClient.Do(req)
			if err == nil {
				_ = response.Body.Close()
			}
		}
		e.WriteAny(socket, NewMessage(msg.OperationId, msg.OperationType, map[string]any{
			"result": "success",
		}))
	default:
		e.WriteAny(socket, NewMessage(msg.OperationId, msg.OperationType, nil))
	}
}

type PortCheckResp struct {
	Code string `json:"code"` // 节点编号
	Use  bool   `json:"use"`  // 是否被占用
	Port string `json:"port"` // 端口
}

func (e *Event) WriteAny(socket *gws.Conn, data any) {
	marshal, _ := json.Marshal(data)
	_ = socket.WriteString(string(marshal))
}
