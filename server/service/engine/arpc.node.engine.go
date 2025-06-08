package engine

import (
	"errors"
	"fmt"
	v1 "github.com/SianHH/frp-package/pkg/config/v1"
	"github.com/lesismal/arpc"
	"server/pkg/utils"
	"server/repository/query"
	"time"
)

func NewARpcNodeEngine(code, ip string, client *arpc.Client) *ARpcNodeEngine {
	return &ARpcNodeEngine{
		code:   code,
		ip:     ip,
		client: client,
	}
}

type ARpcNodeEngine struct {
	code   string
	ip     string
	client *arpc.Client
}

func (e *ARpcNodeEngine) PortCheck(tx *query.Query, ip, port string) error {
	var relay string
	if err := e.client.Call("port_check", port, &relay, time.Second*5); err != nil {
		return err
	}
	if relay != "success" {
		return errors.New(relay)
	}
	return nil
}

func (e *ARpcNodeEngine) Config(tx *query.Query) error {
	node, err := tx.GostNode.Where(tx.GostNode.Code.Eq(e.code)).First()
	if err != nil {
		return err
	}

	serverAddr, serverPort := node.GetAddress()
	var data = ServerConfig{
		Key: node.Code,
		ServerConfig: v1.ServerConfig{
			Auth: v1.AuthServerConfig{
				Token: node.Code,
			},
			BindAddr:                        serverAddr,
			BindPort:                        serverPort,
			KCPBindPort:                     0,
			QUICBindPort:                    0,
			ProxyBindAddr:                   "",
			VhostHTTPPort:                   utils.StrMustInt(node.HttpPort),
			VhostHTTPTimeout:                0,
			VhostHTTPSPort:                  0,
			TCPMuxHTTPConnectPort:           0,
			TCPMuxPassthrough:               false,
			SubDomainHost:                   "",
			Custom404Page:                   "",
			SSHTunnelGateway:                v1.SSHTunnelGateway{},
			WebServer:                       v1.WebServerConfig{},
			EnablePrometheus:                false,
			Log:                             v1.LogConfig{},
			Transport:                       v1.ServerTransportConfig{},
			DetailedErrorsToClient:          nil,
			MaxPortsPerClient:               0,
			UserConnTimeout:                 0,
			UDPPacketSize:                   0,
			NatHoleAnalysisDataReserveHours: 0,
			AllowPorts:                      nil,
			HTTPPlugins: []v1.HTTPPluginOptions{
				{
					Name:      "login-plugins",
					Addr:      "",
					Path:      "/api/v1/public/frp/login",
					Ops:       []string{"Login"},
					TLSVerify: true,
				},
				{
					Name:      "newProxy-plugins",
					Addr:      "",
					Path:      "/api/v1/public/frp/newProxy",
					Ops:       []string{"NewProxy"},
					TLSVerify: true,
				},
				{
					Name:      "closeProxy-plugins",
					Addr:      "",
					Path:      "/api/v1/public/frp/closeProxy",
					Ops:       []string{"CloseProxy"},
					TLSVerify: true,
				},
				{
					Name:      "ping-plugins",
					Addr:      "",
					Path:      "/api/v1/public/frp/ping",
					Ops:       []string{"Ping"},
					TLSVerify: true,
				},
				{
					Name:      "newWorkConn-plugins",
					Addr:      "",
					Path:      "/api/v1/public/frp/newWorkConn",
					Ops:       []string{"NewWorkConn"},
					TLSVerify: true,
				},
				{
					Name:      "newUserConn-plugins",
					Addr:      "",
					Path:      "/api/v1/public/frp/newUserConn",
					Ops:       []string{"NewUserConn"},
					TLSVerify: true,
				},
			},
		},
	}
	//switch node.Protocol {
	//case "quic":
	//	data.ServerConfig.BindPort = 0
	//	data.ServerConfig.QUICBindPort = serverPort
	//case "kcp":
	//	data.ServerConfig.BindPort = 0
	//	data.ServerConfig.KCPBindPort = serverPort
	//}

	var relay string
	if err := e.client.Call("server_config", data, &relay, time.Second*5); err != nil {
		return err
	}
	if relay != "success" {
		return errors.New(relay)
	}
	return nil
}

func (e *ARpcNodeEngine) Ingress(tx *query.Query) error {
	return nil
}

func (e *ARpcNodeEngine) CustomDomain(tx *query.Query, domain, cert, key string, forceHttps int) error {
	node, err := tx.GostNode.Where(tx.GostNode.Code.Eq(e.code)).First()
	if err != nil {
		return err
	}
	var relay string
	if err := e.client.Call("server_domain_config", ServerDomain{
		Domain:     domain,
		Target:     fmt.Sprintf("http://127.0.0.1:%s", node.HttpPort),
		Cert:       cert,
		Key:        key,
		ForceHttps: forceHttps,
	}, &relay, time.Second*5); err != nil {
		return err
	}
	if relay != "success" {
		return errors.New(relay)
	}
	return nil
}

func (e *ARpcNodeEngine) Stop(msg string) {
	_ = e.client.Notify("stop", msg, time.Second*5)
}
