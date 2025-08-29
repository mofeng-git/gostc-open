package bootstrap

import (
	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"
	"os"
	"proxy/configs"
	"proxy/global"
)

func InitServer() {
	_ = os.MkdirAll(global.Config.Certs, 0644)

	cfgBytes, warnMsg, err := global.Config.ParseCaddyFileConfig()
	if err != nil {
		global.Logger.Error("parse caddyfile fail", zap.Error(err))
		Release()
		os.Exit(1)
	}
	for _, item := range warnMsg {
		global.Logger.Warn(item.String())
	}

	config, err := configs.GenerateCaddyConfig(cfgBytes)
	if err != nil {
		global.Logger.Error("generate caddy.Config fail", zap.Error(err))
		Release()
		os.Exit(1)
	}

	if err := caddy.Run(config); err != nil {
		global.Logger.Error("generate caddy.Config fail", zap.Error(err))
		Release()
		os.Exit(1)
	}

	releaseFunc = append(releaseFunc, func() {
		_ = caddy.Stop()
	})
}
