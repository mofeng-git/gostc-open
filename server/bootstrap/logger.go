package bootstrap

import (
	"server/global"
	"server/pkg/logger"
)

func InitLogger() {
	global.Logger = logger.NewLogger(func() logger.Option {
		var to = []string{global.LOGGER_FILE_PATH}
		if global.MODE == "dev" {
			to = append(to, "console")
		}
		return logger.Option{
			To:         to,
			Level:      global.LOGGER_LEVEL,
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
