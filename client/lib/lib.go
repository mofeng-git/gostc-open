package lib

import (
	"crypto/tls"
	"fmt"
	"github.com/go-gost/core/logger"
	"github.com/go-gost/x/config/parsing"
	"github.com/go-gost/x/config/parsing/chain"
	"github.com/go-gost/x/config/parsing/limiter"
	"github.com/go-gost/x/config/parsing/service"
	xlogger "github.com/go-gost/x/logger"
	"github.com/go-gost/x/registry"
	"github.com/lxzan/gws"
	"gostc-sub/common"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var Writer *common.LogWriter

func Init(logLevel string, console bool) {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	level := logger.WarnLevel
	switch logLevel {
	case "trace":
		level = logger.TraceLevel
	case "debug":
		level = logger.DebugLevel
	case "info":
		level = logger.InfoLevel
	case "warn":
		level = logger.WarnLevel
	case "error":
		level = logger.ErrorLevel
	case "fatal":
		level = logger.FatalLevel
	}
	Writer = common.NewLogWriter(console)
	logger.SetDefault(xlogger.NewLogger(xlogger.LevelOption(level), xlogger.OutputOption(Writer)))
	tlsConfig, _ := parsing.BuildDefaultTLSConfig(nil)
	parsing.SetDefaultTLSConfig(tlsConfig)
}

func init() {
	Init("info", false)
}

func GetVersion() string {
	return common.CLIENT_VERSION
}

var svrRunTag = false
var socket *gws.Conn

func RunClient(useTls string, address string, key string) string {
	key = strings.Trim(key, "")
	if key == "" {
		return "请输入Key"
	}
	tlsEnable := useTls == "1"
	svrRunTag = true
	var wsurl = func(tls bool, address string) string {
		var scheme string
		if tlsEnable {
			scheme = "wss"
		} else {
			scheme = "ws"
		}
		return scheme + "://" + address
	}(tlsEnable, address)

	go func() {
		var err error
		for {
			if !svrRunTag {
				return
			}
			socket, _, err = gws.NewClient(common.NewEvent(key, false), &gws.ClientOption{
				Addr:          wsurl + "/api/v1/public/gost/client/ws",
				TlsConfig:     &tls.Config{InsecureSkipVerify: true},
				RequestHeader: http.Header{"key": []string{key}},
				PermessageDeflate: gws.PermessageDeflate{
					Enabled:               true,
					ServerContextTakeover: true,
					ClientContextTakeover: true,
				},
			})
			if err != nil {
				fmt.Println("conn fail,please wait 5 second,retry conn", err)
				time.Sleep(time.Second * 5)
				continue
			}
			_ = socket.WritePing(nil)
			socket.ReadLoop()
		}
	}()
	return "success"
}

func DelClient() {
	// 服务没启动，直接结束
	if !svrRunTag {
		return
	}
	svrRunTag = false
	for k, v := range common.SvcMap {
		if v == true {
			svc := registry.ServiceRegistry().Get(k)
			if svc != nil {
				_ = svc.Close()
			}
		}
	}
	common.SvcMap = make(map[string]bool)
	_ = socket.WriteClose(1000, nil)
}

var TunnelMap = &sync.Map{}

func RunTunnel(useTls string, address string, key string, port string) string {
	tlsEnable := useTls == "1"
	var apiurl = func(tls bool, address string) string {
		var scheme string
		if tlsEnable {
			scheme = "https"
		} else {
			scheme = "http"
		}
		return scheme + "://" + address
	}(tlsEnable, address)

	data, err := common.GetVisitConfig(apiurl + "/api/v1/public/gost/visit?key=" + key)
	if err != nil {
		return "获取配置失败"
	}
	TunnelMap.Store(key, data)
	parseChain, _ := chain.ParseChain(&data.Chain, logger.Default())
	_ = registry.ChainRegistry().Register(data.Chain.Name, parseChain)
	trafficLimiter := limiter.ParseTrafficLimiter(&data.Limiter)
	_ = registry.TrafficLimiterRegistry().Register(data.Limiter.Name, trafficLimiter)
	connLimiter := limiter.ParseConnLimiter(&data.CLimiter)
	_ = registry.ConnLimiterRegistry().Register(data.CLimiter.Name, connLimiter)
	rateLimiter := limiter.ParseRateLimiter(&data.RLimiter)
	_ = registry.RateLimiterRegistry().Register(data.RLimiter.Name, rateLimiter)
	for _, svcCfg := range data.SvcList {
		svcCfg.Addr = ":" + port
		parseService, _ := service.ParseService(&svcCfg)
		go parseService.Serve()
		_ = registry.ServiceRegistry().Register(svcCfg.Name, parseService)
	}
	return "success"
}

func DelTunnel(key string) {
	key = strings.Trim(key, "")
	if key == "" {
		return
	}
	value, ok := TunnelMap.Load(key)
	if !ok {
		return
	}
	cfg := value.(common.VisitCfg)
	for _, svcCfg := range cfg.SvcList {
		_ = registry.ServiceRegistry().Get(svcCfg.Name).Close()
		registry.ServiceRegistry().Unregister(svcCfg.Name)
	}
}
