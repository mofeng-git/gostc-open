package main

import (
	"github.com/go-gost/core/logger"
	xlogger "github.com/go-gost/x/logger"
	"gostc-sub/p2p/frpc"
	"gostc-sub/p2p/pkg/config/types"
	v1 "gostc-sub/p2p/pkg/config/v1"
	log2 "gostc-sub/p2p/pkg/util/log"
	"log"
)

func main() {
	logger.SetDefault(xlogger.NewLogger(xlogger.LevelOption(logger.InfoLevel)))
	log2.RefreshDefault()
	quantity, _ := types.NewBandwidthQuantity("1024KB")
	service := frpc.NewService(v1.ClientCommonConfig{
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

	if err := service.Start(); err != nil {
		log.Fatalln(err)
	}
	service.Wait()
}
