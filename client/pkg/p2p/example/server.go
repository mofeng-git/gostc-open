package main

import (
	"github.com/go-gost/core/logger"
	xlogger "github.com/go-gost/x/logger"
	"gostc-sub/p2p/frps"
	v1 "gostc-sub/p2p/pkg/config/v1"
	log2 "gostc-sub/p2p/pkg/util/log"
	"log"
)

func main() {
	logger.SetDefault(xlogger.NewLogger(xlogger.LevelOption(logger.InfoLevel)))
	log2.RefreshDefault()
	service := frps.NewService(v1.ServerConfig{
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
	if err := service.Start(); err != nil {
		log.Fatalln(err)
	}
	service.Wait()
}
