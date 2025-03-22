package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/kardianos/service"
	"github.com/lxzan/gws"
	"gostc-sub/common"
	"net/http"
	"os"
	"time"
)

var SvcCfg = &service.Config{
	Name:        "gostc",
	DisplayName: "GOSTC",
	Description: "基于GOST开发的内网穿透 客户端/节点",
	Option:      make(service.KeyValue),
}

var Program = &program{}

type program struct {
}

func (p *program) run() {
	var key string
	var tlsEnable bool
	var address string
	var visit bool
	var p2p bool
	var vTunnels string
	var logLevel string
	var version bool
	var console bool
	var server bool
	flag.StringVar(&address, "addr", "gost.sian.one", "server address")
	flag.StringVar(&key, "key", "", "client key")
	flag.BoolVar(&server, "s", false, "server mode")
	flag.BoolVar(&tlsEnable, "tls", true, "enable tls")
	flag.BoolVar(&visit, "v", false, "visit client")
	flag.BoolVar(&p2p, "p2p", false, "p2p client")
	flag.StringVar(&vTunnels, "vts", "", "visit tunnels,example: vkey1:8080,vkey2:8081,vkey3:8082")

	flag.StringVar(&logLevel, "log-level", "error", "log-level trace|debug|info|warn|error|fatal")
	flag.BoolVar(&version, "version", false, "client version")
	flag.BoolVar(&console, "console", false, "log to stdout")
	flag.Parse()

	if version {
		if server {
			fmt.Print(common.SERVER_VERSION)
		} else {
			fmt.Print(common.CLIENT_VERSION)
		}
		os.Exit(0)
	}

	Init(logLevel, console)
	var wsurl = func(tls bool, address string) string {
		var scheme string
		if tlsEnable {
			scheme = "wss"
		} else {
			scheme = "ws"
		}
		return scheme + "://" + address
	}(tlsEnable, address)

	var apiurl = func(tls bool, address string) string {
		var scheme string
		if tlsEnable {
			scheme = "https"
		} else {
			scheme = "http"
		}
		return scheme + "://" + address
	}(tlsEnable, address)

	// 访客模式

	if visit {
		runVisit(vTunnels, apiurl)
		os.Exit(0)
	}

	if p2p {
		runP2P(vTunnels, apiurl)
		os.Exit(0)
	}

	// 连接密钥
	if key == "" {
		fmt.Println("please enter key")
		os.Exit(1)
	}

	fmt.Println("WS_URL：", wsurl)
	fmt.Println("API_URL：", apiurl)

	fullWsUrl := wsurl + "/api/v1/public/gost/client/ws"
	if server {
		fullWsUrl = wsurl + "/api/v1/public/gost/node/ws"
	}

	for {
		socket, _, err := gws.NewClient(common.NewEvent(key, server), &gws.ClientOption{
			Addr:      fullWsUrl,
			TlsConfig: &tls.Config{InsecureSkipVerify: true},
			NewDialer: func() (gws.Dialer, error) {
				return NewDialer(), nil
			},
			RequestHeader: http.Header{"key": []string{key}},
			PermessageDeflate: gws.PermessageDeflate{
				Enabled: true,
			},
		})
		if err != nil {
			fmt.Println("conn fail,please wait 5 second,retry conn", err)
			time.Sleep(time.Second * 5)
			continue
		}
		go func(socket *gws.Conn) {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for {
				select {
				case data := <-Writer.C():
					marshal, _ := json.Marshal(common.NewMessage(uuid.NewString(), "logger", string(data)))
					if err = socket.WriteString(string(marshal)); err != nil {
						fmt.Println("send logger msg fail", err)
						_ = socket.WriteClose(1000, nil)
						return
					}
				case <-ticker.C:
				}
			}
		}(socket)
		socket.ReadLoop()
	}
}

func (p *program) Run() error {
	p.run()
	return nil
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *program) Stop(s service.Service) error {
	return nil
}
