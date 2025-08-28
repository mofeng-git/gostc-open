package main

import (
	"flag"
	"fmt"
	service2 "github.com/SianHH/frp-package/package"
	"github.com/SianHH/frp-package/package/frpc"
	"github.com/SianHH/frp-package/package/frps"
	v1 "github.com/SianHH/frp-package/pkg/config/v1"
	system_service "github.com/kardianos/service"
	"gostc-sub/cli/service_option"
	"gostc-sub/gui/global"
	"gostc-sub/internal/common"
	"gostc-sub/internal/service"
	"gostc-sub/pkg/env"
	"gostc-sub/pkg/signal"
	"gostc-sub/webui/backend/bootstrap"
	"log"
	"os"
	"path/filepath"
	"time"
)

var SvcCfg = &system_service.Config{
	Name:        "gostc",
	DisplayName: "GOSTC",
	Description: "基于FRP开发的内网穿透 客户端/节点",
	Option:      service_option.MakeOptions(),
}

var Program = &program{}

type program struct {
	stopFunc func()
}

func selectMode(cfgFile string, isServer, isVisit, isP2P bool, webAddress string, tunCfg, originType string) string {
	if cfgFile != "" {
		return "cfg"
	}
	if isServer {
		return "server"
	}
	if isVisit {
		return "visit"
	}
	if isP2P {
		return "p2p"
	}
	if webAddress != "" {
		return "ui"
	}
	if tunCfg != "" {
		return "tunCfg"
	}
	if originType == "frps" {
		return "frps"
	}
	if originType == "frpc" {
		return "frpc"
	}
	return "client"
}

func (p *program) run() {
	// 管理端地址
	var address string
	flag.StringVar(&address, "addr", "gost.sian.one", "server address,example: gost.sian.one")
	var tlsEnable bool
	flag.BoolVar(&tlsEnable, "tls", true, "enable tls")

	// 客户端或节点密钥
	var key string
	flag.StringVar(&key, "key", "", "client key")

	// 客户端运行模式
	var server bool
	flag.BoolVar(&server, "s", false, "server mode")
	var visit bool
	flag.BoolVar(&visit, "v", false, "visit client")
	var p2p bool
	flag.BoolVar(&p2p, "p2p", false, "p2p client")
	var wAddress string
	flag.StringVar(&wAddress, "web-addr", "", "web ui address,example: 0.0.0.0:18080")

	var vTunnels string
	flag.StringVar(&vTunnels, "vts", "", "visit tunnels,example: vkey1:8080,vkey2:8081,vkey3:8082")

	var cfgFile string
	flag.StringVar(&cfgFile, "cfg", "", "config file,example: /path/config.yaml")

	// TUN CLI
	var tunCfg string
	flag.StringVar(&tunCfg, "tun-cfg", "", "direct tun by cfg, is internal")
	var tunCfgType string
	flag.StringVar(&tunCfgType, "tun-cfg-type", "", "tun cfg type, node/master, is internal")

	// 其他参数
	var proxyBaseUrl string
	flag.StringVar(&proxyBaseUrl, "proxy-base-url", "", "proxy server api url,example: http://127.0.0.1:8080")
	var version bool
	flag.BoolVar(&version, "version", false, "client version")
	var cfgExample bool
	flag.BoolVar(&cfgExample, "cfg-example", false, "show config example")
	var console bool
	flag.BoolVar(&console, "console", false, "log to stdout")

	// 载入原版配置内容
	var originCfg string
	flag.StringVar(&originCfg, "c", "", "load frp cfg file, /root/frps.toml or /root/frpc.toml")

	flag.Parse()

	address = env.GetEnv("GOSTC_CLIENT_ADDR", address)
	tlsEnable = env.GetEnv("GOSTC_CLIENT_TLS", tlsEnable)

	common.Logger.Console(console)

	if version {
		fmt.Print(common.VERSION)
		os.Exit(0)
	}

	if cfgExample {
		configExample()
		os.Exit(0)
	}

	//var wsurl = common.GenerateWsUrl(tlsEnable, address)
	//var apiurl = common.GenerateHttpUrl(tlsEnable, address)
	generate := common.NewGenerateUrl(tlsEnable, address)

	var mode = selectMode(cfgFile, server, visit, p2p, wAddress, tunCfg, os.Getenv("GOSTC_FRP_TYPE"))

	switch mode {
	case "cfg":
		if err := loadConfig(cfgFile); err != nil {
			log.Fatalln(err)
		}
		startForConfig()
	case "ui":
		basePath, _ := os.Executable()
		basePath = filepath.Dir(basePath)
		moveWebuiCfgDir(basePath, basePath+"/data") // 移动旧配置文件路径
		bootstrap.InitLogger()
		bootstrap.InitFS(basePath)
		bootstrap.InitTodo()
		bootstrap.InitTask()
		bootstrap.InitRouter()
		if err := bootstrap.StartServer(wAddress); err != nil {
			bootstrap.Release()
			fmt.Println(err)
			os.Exit(1)
		}
		p.stopFunc = func() {
			bootstrap.Release()
		}
		<-signal.Free()
	case "visit", "p2p":
		runTunnels(mode, vTunnels, generate)
	case "client", "server":
		if key == "" {
			if mode == "client" {
				fmt.Println("load configuration file", global.BasePath+"/config.yaml")
				if err := loadConfig(global.BasePath + "/config.yaml"); err != nil {
					log.Fatalln(err)
				}
				startForConfig()
				return
			}
			fmt.Println("please enter key")
			os.Exit(1)
		}
		fmt.Println("WS_URL：", generate.WsUrl())
		fmt.Println("API_URL：", generate.HttpUrl())
		var svc service.Service
		switch mode {
		case "server":
			svc = service.NewNode(generate, key, proxyBaseUrl)
		case "client":
			svc = service.NewClient(generate, key)
		}
		if err := svc.Start(); err != nil {
			log.Fatalln("启动失败", err)
		}
		p.stopFunc = func() {
			svc.Stop()
		}
		for {
			if !svc.IsRunning() {
				os.Exit(0)
			}
			time.Sleep(time.Second)
		}
	case "frps":
		file, err := os.ReadFile(originCfg)
		if err != nil {
			fmt.Println("read frps cfg fail", originCfg, err)
			return
		}
		svc, err := frps.NewService(v1.ServerConfig{}, frps.FromBytes(file))
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := svc.Start(); err != nil {
			fmt.Println(err)
			return
		}
		svc.Wait()
	case "frpc":
		file, err := os.ReadFile(originCfg)
		if err != nil {
			fmt.Println("read frpc cfg fail", originCfg, err)
			return
		}
		svc, err := frpc.NewService(v1.ClientCommonConfig{}, nil, nil, frpc.FromBytes(file))
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := svc.Start(); err != nil {
			fmt.Println(err)
			return
		}
		svc.Wait()
	}
}

func (p *program) Run() error {
	p.run()
	return nil
}

func (p *program) Start(s system_service.Service) error {
	go p.run()
	return nil
}
func (p *program) Stop(s system_service.Service) error {
	if p.stopFunc != nil {
		p.stopFunc()
	}
	service2.Range(func(key string, value service2.Service) {
		value.Stop()
	})
	return nil
}

func moveWebuiCfgDir(basePath string, targetPath string) {
	_ = os.MkdirAll(targetPath, 0755)
	if cfgData, err := os.ReadFile(basePath + "/client.json"); err == nil {
		if err = os.WriteFile(targetPath+"/client.json", cfgData, 0644); err == nil {
			_ = os.Remove(basePath + "/client.json")
		}
	}
	if cfgData, err := os.ReadFile(basePath + "/tunnel.json"); err == nil {
		if err = os.WriteFile(targetPath+"/tunnel.json", cfgData, 0644); err == nil {
			_ = os.Remove(basePath + "/tunnel.json")
		}
	}
	if cfgData, err := os.ReadFile(basePath + "/p2p.json"); err == nil {
		if err = os.WriteFile(targetPath+"/p2p.json", cfgData, 0644); err == nil {
			_ = os.Remove(basePath + "/p2p.json")
		}
	}
}
