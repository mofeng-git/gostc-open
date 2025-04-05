package bootstrap

import (
	"go.uber.org/zap"
	"os"
	"proxy/global"
	"proxy/pkg/proxy"
	"time"
)

var server *proxy.Server

func InitServer() {
	var domains = make(map[string]proxy.DomainConfig)
	for domain, item := range global.Config.Domains {
		domains[domain] = proxy.DomainConfig{
			Target: item.Target,
			Cert:   item.Cert,
			Key:    item.Key,
		}
	}

	server = proxy.NewServer(proxy.Config{
		HTTPAddr:  global.Config.HTTPAddr,
		HTTPSAddr: global.Config.HTTPSAddr,
		Default: proxy.DomainConfig{
			Target: global.Config.Default.Target,
			Cert:   global.Config.Default.Cert,
			Key:    global.Config.Default.Key,
		},
		Domains: domains,
	}, global.Logger)

	var err1, err2 error
	go func() {
		err1 = server.StartHTTPServer()
	}()
	go func() {
		err2 = server.StartHTTPSServer()
	}()
	time.Sleep(time.Second)
	releaseFunc = append(releaseFunc, func() {
		server.Close()
	})
	if err1 == nil {
		global.Logger.Info("http server listen on address: " + global.Config.HTTPAddr)
	} else {
		global.Logger.Warn("http server listen on address: "+global.Config.HTTPAddr, zap.Error(err1))
	}
	if err2 == nil {
		global.Logger.Info("https server listen on address: " + global.Config.HTTPSAddr)
	} else {
		global.Logger.Warn("https server listen on address: "+global.Config.HTTPSAddr, zap.Error(err2))
	}
	if err1 != nil && err2 != nil {
		Release()
		global.Logger.Fatal("start server fail")
		os.Exit(1)
	}
}
