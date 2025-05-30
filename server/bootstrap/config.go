package bootstrap

import (
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"server/configs"
	"server/global"
	"server/pkg/env"
	"server/pkg/utils"
	"time"
)

func InitConfig() {
	configFilePath := global.BASE_PATH + "/data/config.yaml"
	_ = os.MkdirAll(filepath.Dir(configFilePath), 0666)

	global.Logger.Info("init config", zap.String("path", configFilePath))
	stat, err := os.Stat(configFilePath)
	if err != nil {
		global.Logger.Warn("stat config fail", zap.String("path", configFilePath), zap.Error(err))
		if err = writeConfigFile(configFilePath); err != nil {
			global.Logger.Fatal("write config fail", zap.String("path", configFilePath), zap.Error(err))
		} else {
			global.Logger.Info("init config finish", zap.String("path", configFilePath))
		}
		return
	}
	if stat.IsDir() {
		global.Logger.Fatal("config is dir", zap.String("path", configFilePath), zap.String("path", configFilePath))
	}
	configFileBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		global.Logger.Fatal("config read fail", zap.String("path", configFilePath), zap.Error(err))
	}
	if err := yaml.Unmarshal(configFileBytes, &global.Config); err != nil {
		global.Logger.Fatal("config serialize fail", zap.String("path", configFilePath), zap.Error(err))
	}
	global.Logger.Info("config load finish", zap.String("path", configFilePath))

	envCfg()
}

// 使用环境变量覆盖一些配置
func envCfg() {
	// 覆盖服务端口
	addr := env.GetEnv("GOSTC_ADMIN_ADDR", "")
	if addr != "" {
		global.Config.Address = addr
		global.Logger.Info("overwrite address configuration with environment variables", zap.String("GOSTC_ADMIN_ADDR", addr))
	}

	fmt.Printf(`
ENV_CONFIG:
GOSTC_ADMIN_ADDR=%s
========================================
`, addr)
}

func writeConfigFile(path string) error {
	global.Config = configs.Config{
		Address:   "0.0.0.0:8080",
		AuthKey:   utils.RandStr(16, utils.AllDict),
		AuthExp:   time.Hour * 24 * 7,
		AuthRenew: time.Hour * 2,
		DbType:    "sqlite",
		Sqlite: configs.Sqlite{
			File:     global.BASE_PATH + "/data/data.db",
			LogLevel: "info",
		},
		Mysql: configs.Mysql{
			Host:     "127.0.0.1",
			Port:     3306,
			DB:       "table_name",
			User:     "root",
			Pwd:      "root",
			Prefix:   "gostc_",
			Extend:   "?timeout=3s&readTimeout=3s&writeTimeout=3s&parseTime=true&loc=Local&charset=utf8mb4,utf8",
			LogLevel: "info",
		},
	}
	marshal, _ := yaml.Marshal(global.Config)
	return os.WriteFile(path, marshal, 0666)
}
