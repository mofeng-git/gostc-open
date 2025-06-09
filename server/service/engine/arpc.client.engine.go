package engine

import (
	"errors"
	"fmt"
	v1 "github.com/SianHH/frp-package/pkg/config/v1"
	"github.com/lesismal/arpc"
	"server/model"
	"server/pkg/utils"
	"server/repository/query"
	"server/service/common/cache"
	"server/service/common/warn_msg"
	"time"
)

type ARpcClientEngine struct {
	code   string
	client *arpc.Client
}

func NewARpcClientEngine(code string, client *arpc.Client) *ARpcClientEngine {
	return &ARpcClientEngine{code: code, client: client}
}

func (e *ARpcClientEngine) Stop(msg string) {
	_ = e.client.Notify("stop", msg, time.Second*5)
}

func (e *ARpcClientEngine) PortCheck(tx *query.Query, ip, port string) error {
	var relay string
	if err := e.client.Call("port_check", port, &relay, time.Second*5); err != nil {
		return err
	}
	if relay != "success" {
		return errors.New(relay)
	}
	return nil
}

func (e *ARpcClientEngine) HostConfig(tx *query.Query, hostCode string) error {
	host, err := tx.GostClientHost.Preload(tx.GostClientHost.Node).Where(tx.GostClientHost.Code.Eq(hostCode)).First()
	if err != nil {
		return err
	}
	warnMsg := warn_msg.GetHostWarnMsg(*host)
	if warnMsg != "" {
		_ = e.RemoveHost(tx, *host, host.Node)
		return errors.New(warnMsg)
	}
	auth, err := tx.GostAuth.Where(tx.GostAuth.TunnelCode.Eq(hostCode)).First()
	if err != nil {
		return err
	}
	var data = HostConfig{
		Key:     host.Code,
		BaseCfg: e.generateServerCommonCfg(host.Node, *auth),
		Http: HTTPProxyConfig{
			HTTPProxyConfig: v1.HTTPProxyConfig{
				ProxyBaseConfig: v1.ProxyBaseConfig{
					Name: host.Code + "_http",
					Type: "http",
					Metadatas: map[string]string{
						"user":     auth.User,
						"password": auth.Password,
					},
					ProxyBackend: v1.ProxyBackend{
						LocalIP:   host.TargetIp,
						LocalPort: utils.StrMustInt(host.TargetPort),
					},
				},
				DomainConfig: v1.DomainConfig{
					CustomDomains: []string{
						host.Node.GetDomainHost(host.DomainPrefix, host.CustomDomain, cache.GetNodeCustomDomain(host.NodeCode)),
					},
				},
				HostHeaderRewrite: host.Node.GetDomainHost(host.DomainPrefix, host.CustomDomain, cache.GetNodeCustomDomain(host.NodeCode)),
			},
			Transport: ProxyTransport{
				UseEncryption:        true,
				UseCompression:       true,
				BandwidthLimit:       fmt.Sprintf("%dKB", host.Limiter*128),
				BandwidthLimitMode:   "server",
				ProxyProtocolVersion: "",
			},
		},
	}
	var relay string
	if err := e.client.Call("host_config", data, &relay, time.Second*5); err != nil {
		return err
	}
	if relay != "success" {
		return errors.New(relay)
	}
	return nil
}

func (e *ARpcClientEngine) RemoveHost(tx *query.Query, host model.GostClientHost, node model.GostNode) error {
	var relay string
	if err := e.client.Call("remove_config", host.Code, &relay, time.Second*5); err != nil {
		return err
	}
	if relay != "success" {
		return errors.New(relay)
	}
	return nil
}

func (e *ARpcClientEngine) ForwardConfig(tx *query.Query, forwardCode string) error {
	forward, err := tx.GostClientForward.Preload(tx.GostClientForward.Node).Where(tx.GostClientForward.Code.Eq(forwardCode)).First()
	if err != nil {
		return err
	}
	warnMsg := warn_msg.GetForwardWarnMsg(*forward)
	if warnMsg != "" {
		_ = e.RemoveForward(tx, *forward)
		return errors.New(warnMsg)
	}

	auth, err := tx.GostAuth.Where(tx.GostAuth.TunnelCode.Eq(forwardCode)).First()
	if err != nil {
		return err
	}

	var data = ForwardConfig{
		Key:     forward.Code,
		BaseCfg: e.generateServerCommonCfg(forward.Node, *auth),
		TCP: TCPProxyConfig{
			TCPProxyConfig: v1.TCPProxyConfig{
				ProxyBaseConfig: v1.ProxyBaseConfig{
					Name: forward.Code + "_tcp",
					Type: "tcp",
					Metadatas: map[string]string{
						"user":     auth.User,
						"password": auth.Password,
					},
					ProxyBackend: v1.ProxyBackend{
						LocalIP:   forward.TargetIp,
						LocalPort: utils.StrMustInt(forward.TargetPort),
					},
				},
				RemotePort: utils.StrMustInt(forward.Port),
			},
			Transport: ProxyTransport{
				UseEncryption:        true,
				UseCompression:       true,
				BandwidthLimit:       fmt.Sprintf("%dKB", forward.Limiter*128),
				BandwidthLimitMode:   "server",
				ProxyProtocolVersion: "",
			},
		},
		UDP: UDPProxyConfig{
			UDPProxyConfig: v1.UDPProxyConfig{
				ProxyBaseConfig: v1.ProxyBaseConfig{
					Name: forward.Code + "_udp",
					Type: "udp",
					Metadatas: map[string]string{
						"user":     auth.User,
						"password": auth.Password,
					},
					ProxyBackend: v1.ProxyBackend{
						LocalIP:   forward.TargetIp,
						LocalPort: utils.StrMustInt(forward.TargetPort),
					},
				},
				RemotePort: utils.StrMustInt(forward.Port),
			},
			Transport: ProxyTransport{
				UseEncryption:        true,
				UseCompression:       true,
				BandwidthLimit:       fmt.Sprintf("%dKB", forward.Limiter*128),
				BandwidthLimitMode:   "server",
				ProxyProtocolVersion: "",
			},
		},
	}

	var relay string
	if err := e.client.Call("forward_config", data, &relay, time.Second*5); err != nil {
		return err
	}
	if relay != "success" {
		return errors.New(relay)
	}
	return nil
}

func (e *ARpcClientEngine) RemoveForward(tx *query.Query, forward model.GostClientForward) error {
	var relay string
	if err := e.client.Call("remove_config", forward.Code, &relay, time.Second*5); err != nil {
		return err
	}
	if relay != "success" {
		return errors.New(relay)
	}
	return nil
}

func (e *ARpcClientEngine) TunnelConfig(tx *query.Query, tunnelCode string) error {
	tunnel, err := tx.GostClientTunnel.Preload(tx.GostClientTunnel.Node).Where(tx.GostClientTunnel.Code.Eq(tunnelCode)).First()
	if err != nil {
		return err
	}
	warnMsg := warn_msg.GetTunnelWarnMsg(*tunnel)
	if warnMsg != "" {
		_ = e.RemoveTunnel(tx, *tunnel, tunnel.Node)
		return errors.New(warnMsg)
	}
	auth, err := tx.GostAuth.Where(tx.GostAuth.TunnelCode.Eq(tunnelCode)).First()
	if err != nil {
		return err
	}
	var data = TunnelConfig{
		Key:     tunnel.Code,
		BaseCfg: e.generateServerCommonCfg(tunnel.Node, *auth),
		STCP: STCPProxyConfig{
			STCPProxyConfig: v1.STCPProxyConfig{
				ProxyBaseConfig: v1.ProxyBaseConfig{
					Name: tunnel.Code + "_stcp",
					Type: "stcp",
					Metadatas: map[string]string{
						"user":     auth.User,
						"password": auth.Password,
					},
					ProxyBackend: v1.ProxyBackend{
						LocalIP:   tunnel.TargetIp,
						LocalPort: utils.StrMustInt(tunnel.TargetPort),
					},
				},
				Secretkey: tunnel.VKey + "_stcp",
			},
			Transport: ProxyTransport{
				UseEncryption:        true,
				UseCompression:       true,
				BandwidthLimit:       fmt.Sprintf("%dKB", tunnel.Limiter*128),
				BandwidthLimitMode:   "server",
				ProxyProtocolVersion: "",
			},
		},
		SUDP: SUDPProxyConfig{
			SUDPProxyConfig: v1.SUDPProxyConfig{
				ProxyBaseConfig: v1.ProxyBaseConfig{
					Name: tunnel.Code + "_sudp",
					Type: "sudp",
					Metadatas: map[string]string{
						"user":     auth.User,
						"password": auth.Password,
					},
					ProxyBackend: v1.ProxyBackend{
						LocalIP:   tunnel.TargetIp,
						LocalPort: utils.StrMustInt(tunnel.TargetPort),
					},
				},
				Secretkey: tunnel.VKey + "_sudp",
			},
			Transport: ProxyTransport{
				UseEncryption:        true,
				UseCompression:       true,
				BandwidthLimit:       fmt.Sprintf("%dKB", tunnel.Limiter*128),
				BandwidthLimitMode:   "server",
				ProxyProtocolVersion: "",
			},
		},
	}

	var relay string
	if err := e.client.Call("tunnel_config", data, &relay, time.Second*5); err != nil {
		return err
	}
	if relay != "success" {
		return errors.New(relay)
	}
	return nil
}

func (e *ARpcClientEngine) RemoveTunnel(tx *query.Query, tunnel model.GostClientTunnel, node model.GostNode) error {
	var relay string
	if err := e.client.Call("remove_config", tunnel.Code, &relay, time.Second*5); err != nil {
		return err
	}
	if relay != "success" {
		return errors.New(relay)
	}
	return nil
}

func (e *ARpcClientEngine) P2PConfig(tx *query.Query, p2pCode string) error {
	p2p, err := tx.GostClientP2P.Preload(tx.GostClientP2P.Node).Where(tx.GostClientP2P.Code.Eq(p2pCode)).First()
	if err != nil {
		return err
	}
	warnMsg := warn_msg.GetP2PWarnMsg(*p2p)
	if warnMsg != "" {
		_ = e.RemoveP2P(tx, *p2p)
		return errors.New(warnMsg)
	}

	auth, err := tx.GostAuth.Where(tx.GostAuth.TunnelCode.Eq(p2p.Code)).First()
	if auth == nil {
		return err
	}

	var data = P2PConfig{
		Key:     p2p.Code,
		BaseCfg: e.generateServerCommonCfg(p2p.Node, *auth),
		XTCP: XTCPProxyConfig{
			XTCPProxyConfig: v1.XTCPProxyConfig{
				ProxyBaseConfig: v1.ProxyBaseConfig{
					Name: p2p.Code + "_p2pxtcp",
					Type: "xtcp",
					Metadatas: map[string]string{
						"user":     auth.User,
						"password": auth.Password,
					},
					ProxyBackend: v1.ProxyBackend{
						LocalIP:   p2p.TargetIp,
						LocalPort: utils.StrMustInt(p2p.TargetPort),
					},
				},
				Secretkey: p2p.VKey,
			},
			Transport: ProxyTransport{
				UseEncryption:  true,
				UseCompression: true,
				//BandwidthLimit:       fmt.Sprintf("%dKB", p2p.Limiter*128),
				ProxyProtocolVersion: "",
			},
		},
		STCP: STCPProxyConfig{
			STCPProxyConfig: v1.STCPProxyConfig{
				ProxyBaseConfig: v1.ProxyBaseConfig{
					Name: p2p.Code + "_p2pstcp",
					Type: "stcp",
					Metadatas: map[string]string{
						"user":     auth.User,
						"password": auth.Password,
					},
					ProxyBackend: v1.ProxyBackend{
						LocalIP:   p2p.TargetIp,
						LocalPort: utils.StrMustInt(p2p.TargetPort),
					},
				},
				Secretkey: p2p.VKey,
			},
			Transport: ProxyTransport{
				UseEncryption:        true,
				UseCompression:       true,
				BandwidthLimit:       fmt.Sprintf("%dKB", p2p.Limiter*128),
				BandwidthLimitMode:   "server",
				ProxyProtocolVersion: "",
			},
		},
	}
	var relay string
	if err := e.client.Call("p2p_config", data, &relay, time.Second*5); err != nil {
		return err
	}
	if relay != "success" {
		return errors.New(relay)
	}
	return nil
}

func (e *ARpcClientEngine) RemoveP2P(tx *query.Query, p2p model.GostClientP2P) error {
	var relay string
	if err := e.client.Call("remove_config", p2p.Code, &relay, time.Second*5); err != nil {
		return err
	}
	if relay != "success" {
		return errors.New(relay)
	}
	return nil
}

func (e *ARpcClientEngine) ProxyConfig(tx *query.Query, proxyCode string) error {
	proxy, err := tx.GostClientProxy.Preload(tx.GostClientProxy.Node).Where(tx.GostClientProxy.Code.Eq(proxyCode)).First()
	if err != nil {
		return err
	}
	warnMsg := warn_msg.GetProxyWarnMsg(*proxy)
	if warnMsg != "" {
		_ = e.RemoveProxy(tx, *proxy)
		return errors.New(warnMsg)
	}

	auth, err := tx.GostAuth.Where(tx.GostAuth.TunnelCode.Eq(proxyCode)).First()
	if err != nil {
		return err
	}

	var data = ProxyConfig{
		Key:      proxy.Code,
		BaseCfg:  e.generateServerCommonCfg(proxy.Node, *auth),
		Name:     proxy.Code + "_proxy",
		Port:     utils.StrMustInt(proxy.Port),
		AuthUser: proxy.AuthUser,
		AuthPwd:  proxy.AuthPwd,
		Metadata: map[string]string{
			"user":     auth.User,
			"password": auth.Password,
		},
		Limiter: fmt.Sprintf("%dKB", proxy.Limiter*128),
	}
	var relay string
	if err := e.client.Call("proxy_config", data, &relay, time.Second*5); err != nil {
		return err
	}
	if relay != "success" {
		return errors.New(relay)
	}
	return nil
}

func (e *ARpcClientEngine) RemoveProxy(tx *query.Query, proxy model.GostClientProxy) error {
	var relay string
	if err := e.client.Call("remove_config", proxy.Code, &relay, time.Second*5); err != nil {
		return err
	}
	if relay != "success" {
		return errors.New(relay)
	}
	return nil
}

func (e *ARpcClientEngine) generateServerCommonCfg(node model.GostNode, auth model.GostAuth) v1.ClientCommonConfig {
	serverAddr, serverPort := node.GetAddress()
	return v1.ClientCommonConfig{
		Auth: v1.AuthClientConfig{
			Token: node.Code,
		},
		ServerAddr: serverAddr,
		ServerPort: serverPort,
		Transport: v1.ClientTransportConfig{
			Protocol: node.Protocol,
		},
		Metadatas: map[string]string{
			"user":     auth.User,
			"password": auth.Password,
		},
		LoginFailExit: new(bool),
	}
}

func (e *ARpcClientEngine) CustomCfgConfig(tx *query.Query, cfgCode string) error {
	cfg, err := tx.FrpClientCfg.Where(tx.FrpClientCfg.Code.Eq(cfgCode)).First()
	if err != nil {
		return err
	}
	warnMsg := warn_msg.GetCfgWarnMsg(*cfg)
	if warnMsg != "" {
		_ = e.RemoveCustomCfg(tx, cfgCode)
		return errors.New(warnMsg)
	}
	var relay string
	if err := e.client.Call("custom_cfg_config", CustomCfgConfig{
		Key:     cfg.Code,
		Type:    cfg.ContentType,
		Content: cfg.Content,
	}, &relay, time.Second*5); err != nil {
		return err
	}
	if relay != "success" {
		return errors.New(relay)
	}
	return nil
}

func (e *ARpcClientEngine) RemoveCustomCfg(tx *query.Query, cfgCode string) error {
	var relay string
	if err := e.client.Call("remove_config", cfgCode, &relay, time.Second*5); err != nil {
		return err
	}
	if relay != "success" {
		return errors.New(relay)
	}
	return nil
}
