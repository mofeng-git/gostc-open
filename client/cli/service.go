package main

import (
	"flag"
	"fmt"
	system_service "github.com/kardianos/service"
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
	Option:      make(system_service.KeyValue),
}

var Program = &program{}

type program struct {
}

func selectMode(isServer, isVisit, isP2P bool, webAddress string) string {
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
	flag.StringVar(&address, "addr", "gost.sian.one", "server address")
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
	flag.StringVar(&wAddress, "web-addr", "", "web ui address")

	var vTunnels string
	flag.StringVar(&vTunnels, "vts", "", "visit tunnels,example: vkey1:8080,vkey2:8081,vkey3:8082")

	// 其他参数
	var proxyBaseUrl string
	flag.StringVar(&proxyBaseUrl, "proxy-base-url", "", "proxy server api url")
	//var logLevel string
	//flag.StringVar(&logLevel, "log-level", "error", "log-level trace|debug|info|warn|error|fatal")
	var version bool
	flag.BoolVar(&version, "version", false, "client version")
	var console bool
	flag.BoolVar(&console, "console", false, "log to stdout")
	flag.Parse()

	common.Logger.Console(console)

	if version {
		fmt.Print(common.VERSION)
		os.Exit(0)
	}

	var wsurl = common.GenerateWsUrl(tlsEnable, address)
	var apiurl = common.GenerateHttpUrl(tlsEnable, address)
	fmt.Println("WS_URL：", wsurl)
	fmt.Println("API_URL：", apiurl)

	var mode = selectMode(server, visit, p2p, wAddress)

	switch mode {
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
	case "client", "server":
		if key == "" {
			fmt.Println("please enter key")
			os.Exit(1)
		}
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
