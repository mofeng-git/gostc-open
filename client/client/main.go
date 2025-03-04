package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/go-gost/core/logger"
	"github.com/go-gost/x/config/parsing"
	"github.com/go-gost/x/config/parsing/chain"
	"github.com/go-gost/x/config/parsing/limiter"
	"github.com/go-gost/x/config/parsing/service"
	xlogger "github.com/go-gost/x/logger"
	"github.com/go-gost/x/registry"
	"github.com/google/uuid"
	"github.com/lxzan/gws"
	"gostc-sub/common"
	"gostc-sub/pkg/signal"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var Writer *common.LogWriter

func Init(logLevel string, console bool) {
	fixDNSResolver()
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	level := logger.InfoLevel
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

func fixDNSResolver() {
	if net.DefaultResolver != nil {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err := net.DefaultResolver.LookupHost(timeoutCtx, "google.com")
		if err == nil {
			return
		}
	}
	net.DefaultResolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if addr == "127.0.0.1:53" || addr == "[::1]:53" {
				addr = "8.8.8.8:53"
			}
			var d net.Dialer
			return d.DialContext(ctx, network, addr)
		},
	}
}

type MyDialer struct {
	net.Resolver
	ctx context.Context
}

func NewDialer() *MyDialer {
	return &MyDialer{
		Resolver: net.Resolver{
			PreferGo:     true,
			StrictErrors: false,
			Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
				var split = strings.Split(addr, ":")
				if len(split) != 2 {
					return nil, errors.New("addr analysis fail")
				}
				host, err := net.LookupHost(split[0])
				if len(host) == 0 {
					return nil, err
				}
				var d net.Dialer
				return d.DialContext(ctx, network, host[0]+":"+split[1])
			},
		},
		ctx: context.Background(),
	}
}

func (m *MyDialer) Dial(network, addr string) (c net.Conn, err error) {
	return m.Resolver.Dial(m.ctx, network, addr)
}

func main() {
	var key string
	var tlsEnable bool
	var address string
	var visit bool
	var vTunnels string
	var logLevel string
	var version bool
	var console bool
	var server bool
	flag.StringVar(&address, "addr", "gost.sian.one", "server address")
	flag.StringVar(&key, "key", "", "client key")
	flag.BoolVar(&server, "s", false, "server mode")
	flag.BoolVar(&tlsEnable, "tls", true, "enable tls")
	flag.BoolVar(&visit, "v", false, "visit client")
	flag.StringVar(&vTunnels, "vts", "", "visit tunnels,example: vkey1:8080,vkey2:8081,vkey3:8082")

	flag.StringVar(&logLevel, "log-level", "error", "log-level trace|debug|info|warn|error|fatal")
	flag.BoolVar(&version, "version", false, "client version")
	flag.BoolVar(&console, "console", false, "log to stdout")
	flag.Parse()

	if version {
		if server {
			fmt.Print(common.SERVER_VERSION)
		} else {
			fmt.Print(common.CLIENT_VERSION)
		}
		os.Exit(0)
	}

	Init(logLevel, console)
	var wsurl = func(tls bool, address string) string {
		var scheme string
		if tlsEnable {
			scheme = "wss"
		} else {
			scheme = "ws"
		}
		return scheme + "://" + address
	}(tlsEnable, address)

	var apiurl = func(tls bool, address string) string {
		var scheme string
		if tlsEnable {
			scheme = "https"
		} else {
			scheme = "http"
		}
		return scheme + "://" + address
	}(tlsEnable, address)

	// 访客模式

	if visit {
		runVisit(vTunnels, apiurl)
		os.Exit(0)
	}

	// 连接密钥
	if key == "" {
		fmt.Println("please enter key")
		return
	}

	fmt.Println("WS_URL：", wsurl)
	fmt.Println("API_URL：", apiurl)

	fullWsUrl := wsurl + "/api/v1/public/gost/client/ws"
	if server {
		fullWsUrl = wsurl + "/api/v1/public/gost/node/ws"
	}

	for {
		socket, _, err := gws.NewClient(common.NewEvent(key, server), &gws.ClientOption{
			Addr:      fullWsUrl,
			TlsConfig: &tls.Config{InsecureSkipVerify: true},
			NewDialer: func() (gws.Dialer, error) {
				return NewDialer(), nil
			},
			RequestHeader: http.Header{"key": []string{key}},
			PermessageDeflate: gws.PermessageDeflate{
				Enabled: true,
			},
		})
		if err != nil {
			fmt.Println("conn fail,please wait 5 second,retry conn", err)
			time.Sleep(time.Second * 5)
			continue
		}
		go func(socket *gws.Conn) {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for {
				select {
				case data := <-Writer.C():
					marshal, _ := json.Marshal(common.NewMessage(uuid.NewString(), "logger", string(data)))
					if err = socket.WriteString(string(marshal)); err != nil {
						fmt.Println("send logger msg fail", err)
						_ = socket.WriteClose(1000, nil)
						return
					}
				case <-ticker.C:
				}
			}
		}(socket)
		socket.ReadLoop()
	}
}

type Tunnel struct {
	VKey string
	Port string
}

func runVisit(vTunnels string, apiurl string) {
	var tunnelList []Tunnel
	tunnels := strings.Split(vTunnels, ",")
	for _, tunnel := range tunnels {
		kp := strings.Split(tunnel, ":")
		if len(kp) != 2 {
			continue
		}
		tunnelList = append(tunnelList, Tunnel{
			VKey: kp[0],
			Port: kp[1],
		})
	}

	for _, tunnel := range tunnelList {
		data, err := common.GetVisitConfig(apiurl + "/api/v1/public/gost/visit?key=" + tunnel.VKey)
		if err != nil {
			fmt.Println("获取隧道配置失败", tunnel.VKey, tunnel.Port, err)
			continue
		}
		parseChain, _ := chain.ParseChain(&data.Chain, logger.Default())
		_ = registry.ChainRegistry().Register(data.Chain.Name, parseChain)
		trafficLimiter := limiter.ParseTrafficLimiter(&data.Limiter)
		_ = registry.TrafficLimiterRegistry().Register(data.Limiter.Name, trafficLimiter)
		connLimiter := limiter.ParseConnLimiter(&data.CLimiter)
		_ = registry.ConnLimiterRegistry().Register(data.CLimiter.Name, connLimiter)
		rateLimiter := limiter.ParseRateLimiter(&data.RLimiter)
		_ = registry.RateLimiterRegistry().Register(data.RLimiter.Name, rateLimiter)
		for _, svcCfg := range data.SvcList {
			svcCfg.Addr = ":" + tunnel.Port
			parseService, _ := service.ParseService(&svcCfg)
			go parseService.Serve()
			_ = registry.ServiceRegistry().Register(svcCfg.Name, parseService)
		}
		fmt.Println("隧道启动成功", tunnel.VKey, "0.0.0.0:"+tunnel.Port)
	}
	<-signal.Free()
}
