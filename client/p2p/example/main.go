package main

import (
	"fmt"
	"github.com/go-gost/core/logger"
	xlogger "github.com/go-gost/x/logger"
	"gostc-sub/p2p/frpc"
	"gostc-sub/p2p/frps"
	"gostc-sub/p2p/pkg/config/types"
	v1 "gostc-sub/p2p/pkg/config/v1"
	log2 "gostc-sub/p2p/pkg/util/log"
	"gostc-sub/p2p/registry"
	"log"
	"time"
)

func main() {
	logger.SetDefault(xlogger.NewLogger(xlogger.LevelOption(logger.InfoLevel)))
	log2.RefreshDefault()

	server := NewServer()
	if err := server.Start(); err != nil {
		log.Fatalln(err)
	}
	fmt.Println(registry.Set("server", server))

	client := NewClient()
	if err := client.Start(); err != nil {
		log.Fatalln(err)
	}
	fmt.Println(registry.Set("client", client))

	visit := NewVisit()
	if err := visit.Start(); err != nil {
		log.Fatalln(err)
	}
	fmt.Println(registry.Set("visit", visit))

	time.Sleep(time.Hour)
}

func NewServer() *frps.Service {
	return frps.NewService(v1.ServerConfig{
		Auth:                            v1.AuthServerConfig{},
		BindAddr:                        "",
		BindPort:                        0,
		Transport:                       v1.ServerTransportConfig{},
		DetailedErrorsToClient:          nil,
		MaxPortsPerClient:               0,
		UserConnTimeout:                 0,
		NatHoleAnalysisDataReserveHours: 0,
		AllowPorts:                      nil,
		HTTPPlugins: []v1.HTTPPluginOptions{
			{
				Name:      "login-plugin",
				Addr:      "http://127.0.0.1:8080",
				Path:      "/api/v1/public/p2p/login",
				Ops:       []string{"Login"},
				TLSVerify: false,
			},
			{
				Name:      "new-plugin",
				Addr:      "http://127.0.0.1:8080",
				Path:      "/api/v1/public/p2p/new",
				Ops:       []string{"NewProxy"},
				TLSVerify: false,
			},
		},
	})
}

func NewClient() *frpc.Service {
	quantity, _ := types.NewBandwidthQuantity("128KB")
	return frpc.NewService(v1.ClientCommonConfig{
		Auth:              v1.AuthClientConfig{},
		User:              "user1",
		ServerAddr:        "",
		ServerPort:        0,
		NatHoleSTUNServer: "",
		DNSServer:         "",
		LoginFailExit:     nil,
		Transport:         v1.ClientTransportConfig{},
		Metadatas:         nil,
	}, []v1.ProxyConfigurer{
		&v1.STCPProxyConfig{
			ProxyBaseConfig: v1.ProxyBaseConfig{
				Name: "test-stcp",
				Type: "stcp",
				Transport: v1.ProxyTransport{
					UseEncryption:        true,
					UseCompression:       true,
					BandwidthLimit:       quantity,
					BandwidthLimitMode:   "client",
					ProxyProtocolVersion: "",
				},
				ProxyBackend: v1.ProxyBackend{
					LocalIP:   "192.168.0.172",
					LocalPort: 22714,
				},
			},
			Secretkey:  "stcp-secret",
			AllowUsers: []string{"*"},
		},
		&v1.XTCPProxyConfig{
			ProxyBaseConfig: v1.ProxyBaseConfig{
				Name: "test-xtcp",
				Type: "xtcp",
				ProxyBackend: v1.ProxyBackend{
					LocalIP:   "192.168.0.172",
					LocalPort: 22714,
				},
			},
			Secretkey:  "xtcp-secret",
			AllowUsers: []string{"*"},
		},
	}, []v1.VisitorConfigurer{})
}

func NewVisit() *frpc.Service {
	return frpc.NewService(v1.ClientCommonConfig{
		Auth:              v1.AuthClientConfig{},
		User:              "",
		ServerAddr:        "",
		ServerPort:        0,
		NatHoleSTUNServer: "",
		DNSServer:         "",
		LoginFailExit:     nil,
		Transport:         v1.ClientTransportConfig{},
		Metadatas:         nil,
	}, []v1.ProxyConfigurer{}, []v1.VisitorConfigurer{
		&v1.STCPVisitorConfig{
			VisitorBaseConfig: v1.VisitorBaseConfig{
				Name:       "test-stcp",
				Type:       "stcp",
				ServerName: "test-stcp",
				SecretKey:  "stcp-secret",
				BindAddr:   "127.0.0.1",
				//BindPort:   0, // 0 不监听端口，仅供xtcp使用
				BindPort: 16000,
			},
		},
		//&v1.XTCPVisitorConfig{
		//	VisitorBaseConfig: v1.VisitorBaseConfig{
		//		Name:       "test-xtcp",
		//		Type:       "xtcp",
		//		Transport:  v1.VisitorTransport{},
		//		SecretKey:  "xtcp-secret",
		//		ServerName: "test-xtcp",
		//		BindAddr:   "",
		//		BindPort:   16000,
		//	},
		//	Protocol:          "",   // 隧道底层通信协议，可选 quic 和 kcp，默认为 quic。
		//	KeepTunnelOpen:    true, // 是否保持隧道打开，如果开启，会定期检查隧道状态并尝试保持打开。
		//	MaxRetriesAnHour:  8,    // 每小时最大重试次数8次
		//	MinRetryInterval:  90,   // 最小重试间隔 90秒
		//	FallbackTo:        "test-stcp",
		//	FallbackTimeoutMs: 5000,
		//},
	})
}
