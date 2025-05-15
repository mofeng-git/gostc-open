package main

import (
	"flag"
	"fmt"
	system_service "github.com/kardianos/service"
	"gostc-sub/cli/service_option"
	"gostc-sub/gui/global"
	"gostc-sub/internal/common"
	"gostc-sub/internal/service"
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
	Description: "基于GOST开发的内网穿透 客户端/节点",
	Option:      service_option.MakeOptions(),
}

var Program = &program{}

type program struct {
}

func selectMode(cfgFile string, isServer, isVisit, isP2P bool, webAddress string, tun bool) string {
	if tun {
		return "tun"
	}
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
	var tun bool
	flag.BoolVar(&tun, "tun", false, "tun client")
	var wAddress string
	flag.StringVar(&wAddress, "web-addr", "", "web ui address,example: 0.0.0.0:18080")

	var vTunnels string
	flag.StringVar(&vTunnels, "vts", "", "visit tunnels,example: vkey1:8080,vkey2:8081,vkey3:8082")

	var cfgFile string
	flag.StringVar(&cfgFile, "cfg", "", "config file,example: /path/config.yaml")

	// 其他参数
	var proxyBaseUrl string
	flag.StringVar(&proxyBaseUrl, "proxy-base-url", "", "proxy server api url,example: http://127.0.0.1:8080")
	var version bool
	flag.BoolVar(&version, "version", false, "client version")
	var cfgExample bool
	flag.BoolVar(&cfgExample, "cfg-example", false, "show config example")
	var console bool
	flag.BoolVar(&console, "console", false, "log to stdout")

	flag.Parse()

	common.Logger.Console(console)

	if version {
		fmt.Print(common.VERSION)
		os.Exit(0)
	}

	if cfgExample {
		configExample()
		os.Exit(0)
	}

	var wsurl = common.GenerateWsUrl(tlsEnable, address)
	var apiurl = common.GenerateHttpUrl(tlsEnable, address)

	var mode = selectMode(cfgFile, server, visit, p2p, wAddress, tun)

	switch mode {
	case "cfg":
		if err := loadConfig(cfgFile); err != nil {
			log.Fatalln(err)
		}
		startForConfig()
	case "ui":
		basePath, _ := os.Executable()
		basePath = filepath.Dir(basePath)
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
		<-signal.Free()
		bootstrap.Release()
	case "visit", "p2p":
		runTunnels(mode, vTunnels, apiurl)
	case "client", "server", "tun":
		if key == "" {
			if mode == "client" {
				fmt.Println("load configuration file", global.BasePath+"/config.yaml")
				if err := loadConfig(cfgFile); err != nil {
					log.Fatalln(err)
				}
				startForConfig()
				return
			}
			fmt.Println("please enter key")
			os.Exit(1)
		}
		fmt.Println("WS_URL：", wsurl)
		fmt.Println("API_URL：", apiurl)
		var svc service.Service
		switch mode {
		case "server":
			svc = service.NewServer(wsurl, key, proxyBaseUrl)
		case "client":
			svc = service.NewClient(wsurl, key)
		}
		if err := svc.Start(); err != nil {
			log.Fatalln("启动失败", err)
		}
		for {
			if !svc.IsRunning() {
				os.Exit(0)
			}
			time.Sleep(time.Second)
		}
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
	return nil
}
