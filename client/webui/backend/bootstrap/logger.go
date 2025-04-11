package bootstrap

import (
	"gostc-sub/pkg/logger"
	"gostc-sub/webui/backend/global"
)

func InitLogger() {
	global.Logger = logger.NewLogger(func() logger.Option {
		return logger.Option{
			To:         []string{"console"},
			Level:      "info",
			MaxSize:    100,
			MaxAge:     30,
			MaxBackups: 10,
			Compress:   true,
		}
	}())

	releaseFunc = append(releaseFunc, func() {
		_ = global.Logger.Sync()
	})
}
