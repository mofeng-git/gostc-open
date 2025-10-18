package bootstrap

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/radovskyb/watcher"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"proxy/configs"
	"proxy/global"
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

	go watchConfigReload(configFilePath)
}

func watchConfigReload(configFilePath string) {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write)

	delay := 3 * time.Second // 延迟更新时间，避免反复更新导致配置出现遗漏
	var delayTimer *time.Timer

	go func() {
		for {
			select {
			case <-w.Event:
				global.Logger.Info("watch config reload", zap.String("path", configFilePath))
				if delayTimer != nil {
					delayTimer.Stop()
				}
				delayTimer = time.AfterFunc(delay, func() {
					reloadCaddyConfig(configFilePath)
				})
			case _ = <-w.Error:
				return
			case <-w.Closed:
				return
			}
		}
	}()
	if err := w.Add(configFilePath); err != nil {
		global.Logger.Warn("watcher config file failed", zap.Error(err))
		return
	}
	if err := w.Start(time.Second); err != nil {
		global.Logger.Warn("watcher config file failed", zap.Error(err))
		return
	}
	global.Logger.Info("watcher config file success")
}

func reloadCaddyConfig(configFilePath string) {
	configFileBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		global.Logger.Error("reload config fail", zap.String("path", configFilePath), zap.Error(err))
		return
	}
	var newConfig = configs.Config{}
	if err := yaml.Unmarshal(configFileBytes, &newConfig); err != nil {
		global.Logger.Error("reload config fail", zap.String("path", configFilePath), zap.Error(err))
		return
	}

	cfgBytes, _, err := newConfig.ParseCaddyFileConfig()
	if err != nil {
		global.Logger.Error("reload config fail, parse caddyfile fail", zap.String("path", configFilePath), zap.Error(err))
		return
	}

	if err := caddy.Load(cfgBytes, false); err != nil {
		global.Logger.Error("reload config fail, reload caddyfile fail", zap.String("path", configFilePath), zap.Error(err))
		return
	}

	global.Logger.Info("reload caddyfile", zap.String("caddyfile", string(cfgBytes)))
	global.Config = &newConfig
}

func writeConfigFile(path string) error {
	global.Config = &configs.Config{
		HTTPAddr:  "0.0.0.0:80",
		HTTPSAddr: "0.0.0.0:443",
		Certs:     global.BASE_PATH + "/data/certs",
		ApiAddr:   "0.0.0.0:8080",
		Default: configs.DomainConfig{
			Target: "http://127.0.0.1:8080",
			Cert:   global.BASE_PATH + "/data/certs/default.pem",
			Key:    global.BASE_PATH + "/data/certs/default.key",
		},
		Domains: map[string]configs.DomainConfig{
			"www.example.com": {
				Target: "http://127.0.0.1:8080",
				Cert:   "",
				Key:    "",
			},
		},
	}
	marshal, _ := yaml.Marshal(global.Config)
	return os.WriteFile(path, marshal, 0666)
}
