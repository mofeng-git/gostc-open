package main

import (
	"github.com/go-gost/core/logger"
	xlogger "github.com/go-gost/x/logger"
	"gostc-sub/p2p/frpc"
	v1 "gostc-sub/p2p/pkg/config/v1"
	log2 "gostc-sub/p2p/pkg/util/log"
	"log"
)

func main() {
	logger.SetDefault(xlogger.NewLogger(xlogger.LevelOption(logger.InfoLevel)))
	log2.RefreshDefault()
	service := frpc.NewService(v1.ClientCommonConfig{
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
				BindPort:   0, // 0 不监听端口，仅供xtcp使用
				//BindPort: 16000,
			},
		},
		&v1.XTCPVisitorConfig{
			VisitorBaseConfig: v1.VisitorBaseConfig{
				Name:       "test-xtcp",
				Type:       "xtcp",
				Transport:  v1.VisitorTransport{},
				SecretKey:  "xtcp-secret",
				ServerName: "test-xtcp",
				BindAddr:   "",
				BindPort:   16000,
			},
			Protocol:          "",   // 隧道底层通信协议，可选 quic 和 kcp，默认为 quic。
			KeepTunnelOpen:    true, // 是否保持隧道打开，如果开启，会定期检查隧道状态并尝试保持打开。
			MaxRetriesAnHour:  8,    // 每小时最大重试次数8次
			MinRetryInterval:  90,   // 最小重试间隔 90秒
			FallbackTo:        "test-stcp",
			FallbackTimeoutMs: 5000,
		},
	})

	if err := service.Start(); err != nil {
		log.Fatalln(err)
	}
	service.Wait()
}
