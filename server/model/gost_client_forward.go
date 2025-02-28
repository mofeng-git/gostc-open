package model

import (
	"encoding/json"
	"fmt"
	"github.com/go-gost/x/config"
	"server/pkg/utils"
	"strings"
)

type GostClientForward struct {
	Base
	Name          string     `gorm:"column:name;index;comment:名称"`
	TargetIp      string     `gorm:"column:target_ip;index;comment:内网IP"`
	TargetPort    string     `gorm:"column:target_port;index;comment:内网端口"`
	Port          string     `gorm:"column:port;comment:访问端口"`
	NoDelay       int        `gorm:"column:no_delay;size:1;comment:无等待延迟"`
	NodeCode      string     `gorm:"column:node_code;index;comment:节点编号"`
	Node          GostNode   `gorm:"foreignKey:NodeCode;references:Code"`
	ClientCode    string     `gorm:"column:client_code;index;comment:客户端编号"`
	Client        GostClient `gorm:"foreignKey:ClientCode;references:Code"`
	UserCode      string     `gorm:"column:user_code;index;comment:用户编号"`
	User          SystemUser `gorm:"foreignKey:UserCode;references:Code"`
	Enable        int        `gorm:"column:enable;size:1;default:1;comment:启用状态"`
	Status        int        `gorm:"column:status;size:1;default:1;comment:状态"`
	MatcherEnable int        `gorm:"column:matcher_enable;size:1;default:2;comment:是否开启匹配规则"`
	Matcher       string     `gorm:"column:matcher;comment:匹配规则"`
	TcpMatcher    string     `gorm:"column:tcp_matcher;comment:规则匹配"`
	SSHMatcher    string     `gorm:"column:ssh_matcher;comment:规则匹配"`
	GostClientAdmission
	GostClientConfig
}

func (forward *GostClientForward) GetTcpMatcher() (targetIp, targetPort string) {
	split := strings.Split(forward.TcpMatcher, "$$")
	if len(split) == 2 {
		return split[0], split[1]
	}
	return "", ""
}

func (forward *GostClientForward) SetTcpMatcher(targetIp, targetPort string) {
	if !utils.ValidateLocalIP(targetIp) {
		return
	}
	if !utils.ValidatePort(targetPort) {
		return
	}
	forward.TcpMatcher = targetIp + "$$" + targetPort
}

func (forward *GostClientForward) GetSSHMatcher() (targetIp, targetPort string) {
	split := strings.Split(forward.SSHMatcher, "$$")
	if len(split) == 2 {
		return split[0], split[1]
	}
	return "", ""
}
func (forward *GostClientForward) SetSSHMatcher(targetIp, targetPort string) {
	if !utils.ValidateLocalIP(targetIp) {
		return
	}
	if !utils.ValidatePort(targetPort) {
		return
	}
	forward.SSHMatcher = targetIp + "$$" + targetPort
}

func (forward *GostClientForward) GetMatcher() (result []ForwardMatcher) {
	_ = json.Unmarshal([]byte(forward.Matcher), &result)
	return result
}

func (forward *GostClientForward) SetMatcher(list []ForwardMatcher) {
	var validMatcher []ForwardMatcher
	for _, item := range list {
		host := strings.ReplaceAll(item.Host, " ", "")
		targetIp := strings.ReplaceAll(item.TargetIp, " ", "")
		targetPort := strings.ReplaceAll(item.TargetPort, " ", "")
		if host == "" {
			continue
		}
		if !utils.ValidateLocalIP(targetIp) {
			continue
		}
		if !utils.ValidatePort(targetPort) {
			continue
		}
		validMatcher = append(validMatcher, ForwardMatcher{
			Host:       host,
			TargetIp:   targetIp,
			TargetPort: targetPort,
		})
	}
	marshal, _ := json.Marshal(validMatcher)
	forward.Matcher = string(marshal)
}

type ForwardMatcher struct {
	Host       string `json:"host"`
	TargetIp   string `json:"targetIp"`
	TargetPort string `json:"targetPort"`
}

func (forward *GostClientForward) GenerateTcpSvcConfig(chain, limiter, cLimiter, rLimiter, obs, admissionWhite, admissionBlack string) (clientCfg config.ServiceConfig, ok bool) {
	if forward.Node.Forward != 1 {
		return clientCfg, ok
	}
	var target = forward.TargetIp + ":" + forward.TargetPort

	var forwardNodes []*config.ForwardNodeConfig
	if forward.MatcherEnable == 1 {
		for _, matcher := range forward.GetMatcher() {
			forwardNodes = append(forwardNodes, &config.ForwardNodeConfig{
				Name: matcher.TargetIp + ":" + matcher.TargetPort,
				Addr: matcher.TargetIp + ":" + matcher.TargetPort,
				Matcher: &config.NodeMatcherConfig{
					Rule: fmt.Sprintf("Host(`%s`)", matcher.Host),
				},
			})
		}
		sshIp, sshPort := forward.GetSSHMatcher()
		if sshIp != "" && sshPort != "" {
			var addr = sshIp + ":" + sshPort
			forwardNodes = append(forwardNodes, &config.ForwardNodeConfig{
				Name: addr,
				Addr: addr,
				Matcher: &config.NodeMatcherConfig{
					Rule: "Proto(`ssh`)",
				},
			})
		}
		tcpIp, tcpPort := forward.GetTcpMatcher()
		if tcpIp != "" && tcpPort != "" {
			var addr = tcpIp + ":" + tcpPort
			forwardNodes = append(forwardNodes, &config.ForwardNodeConfig{
				Name: addr,
				Addr: addr,
			})
		}
	} else {
		forwardNodes = append(forwardNodes, &config.ForwardNodeConfig{
			Name: target,
			Addr: target,
		})
	}

	var admissions []string
	if forward.WhiteEnable == 1 {
		admissions = append(admissions, admissionWhite)
	}
	if forward.BlackEnable == 1 {
		admissions = append(admissions, admissionBlack)
	}
	clientCfg = config.ServiceConfig{
		Name:       "tcp_" + forward.Code,
		Addr:       ":" + forward.Port,
		Admissions: admissions,
		Limiter:    limiter,
		CLimiter:   cLimiter,
		RLimiter:   rLimiter,
		Observer:   obs,
		Recorders:  nil,
		Handler:    &config.HandlerConfig{Type: "rtcp", Metadata: map[string]any{"keepAlive": true, "sniffing": true, "nodelay": true}},
		Listener:   &config.ListenerConfig{Type: "rtcp", Chain: chain, Metadata: map[string]any{"keepAlive": true, "nodelay": true}},
		Forwarder: &config.ForwarderConfig{
			Nodes: forwardNodes,
		},
		Metadata: map[string]any{
			"keepAlive":             true,
			"enableStats":           true,
			"observer.period":       "60s",
			"observer.resetTraffic": true,
			"nodelay":               true,
		},
	}
	return clientCfg, true
}

func (forward *GostClientForward) GenerateUdpSvcConfig(chain, limiter, cLimiter, rLimiter, obs, admissionWhite, admissionBlack string) (config.ServiceConfig, bool) {
	if forward.Node.Forward != 1 {
		return config.ServiceConfig{}, false
	}
	var target = forward.TargetIp + ":" + forward.TargetPort

	var forwardNodes []*config.ForwardNodeConfig
	if forward.MatcherEnable == 1 {
		for _, matcher := range forward.GetMatcher() {
			forwardNodes = append(forwardNodes, &config.ForwardNodeConfig{
				Name: matcher.TargetIp + ":" + matcher.TargetPort,
				Addr: matcher.TargetIp + ":" + matcher.TargetPort,
				Matcher: &config.NodeMatcherConfig{
					Rule: fmt.Sprintf("Host(`%s`)", matcher.Host),
				},
			})
		}
		sshIp, sshPort := forward.GetSSHMatcher()
		if sshIp != "" && sshPort != "" {
			var addr = sshIp + ":" + sshPort
			forwardNodes = append(forwardNodes, &config.ForwardNodeConfig{
				Name: addr,
				Addr: addr,
				Matcher: &config.NodeMatcherConfig{
					Rule: "Proto(`ssh`)",
				},
			})
		}
		tcpIp, tcpPort := forward.GetTcpMatcher()
		if tcpIp != "" && tcpPort != "" {
			var addr = tcpIp + ":" + tcpPort
			forwardNodes = append(forwardNodes, &config.ForwardNodeConfig{
				Name: addr,
				Addr: addr,
			})
		}
	} else {
		forwardNodes = append(forwardNodes, &config.ForwardNodeConfig{
			Name: target,
			Addr: target,
		})
	}

	var admissions []string
	if forward.WhiteEnable == 1 {
		admissions = append(admissions, admissionWhite)
	}
	if forward.BlackEnable == 1 {
		admissions = append(admissions, admissionBlack)
	}
	return config.ServiceConfig{
		Name:       "udp_" + forward.Code,
		Addr:       ":" + forward.Port,
		Admissions: admissions,
		Limiter:    limiter,
		CLimiter:   cLimiter,
		RLimiter:   rLimiter,
		Observer:   obs,
		Handler:    &config.HandlerConfig{Type: "rudp", Metadata: map[string]any{"keepAlive": true, "sniffing": true, "nodelay": true}},
		Listener:   &config.ListenerConfig{Type: "rudp", Chain: chain, Metadata: map[string]any{"keepAlive": true, "nodelay": true}},
		Forwarder: &config.ForwarderConfig{
			Nodes: forwardNodes,
		},
		Metadata: map[string]any{
			"keepAlive":             true,
			"enableStats":           true,
			"observer.period":       "60s",
			"observer.resetTraffic": true,
			"nodelay":               true,
		},
	}, true
}

func (forward *GostClientForward) GenerateChainConfig(auth GostAuth) config.ChainConfig {
	var metadata = make(map[string]any)
	_ = json.Unmarshal([]byte(forward.Node.ForwardMetadata), &metadata)
	if forward.NoDelay == 1 {
		metadata["nodelay"] = true
	}
	metadata["keepAlive"] = true

	var protocol, address string
	protocol = forward.Node.Protocol
	address = forward.Node.Address + ":" + forward.Node.ForwardConnPort
	replaceAddress := strings.Split(forward.Node.ForwardReplaceAddress, "://")
	if len(replaceAddress) == 2 {
		protocol = replaceAddress[0]
		address = replaceAddress[1]
	}

	return config.ChainConfig{
		Name: "chain_" + forward.Code,
		Hops: []*config.HopConfig{
			{
				Nodes: []*config.NodeConfig{
					{
						Addr: address,
						Connector: &config.ConnectorConfig{
							Type:     "relay",
							Metadata: metadata,
							Auth: &config.AuthConfig{
								Username: auth.User,
								Password: auth.Password,
							},
						},
						Dialer: &config.DialerConfig{
							Type:     protocol,
							Metadata: metadata,
						},
					},
				},
			},
		},
	}
}

func (forward *GostClientForward) GenerateLimiter() config.LimiterConfig {
	return config.LimiterConfig{
		Name: "limiter_" + forward.Code,
		Limits: []string{
			fmt.Sprintf("$ %dKB  %dKB", forward.Limiter*128, forward.Limiter*128),
		},
	}
}

func (forward *GostClientForward) GenerateCLimiter() config.LimiterConfig {
	return config.LimiterConfig{
		Name: "cLimiter_" + forward.Code,
		Limits: []string{
			fmt.Sprintf("$ %d", forward.CLimiter),
		},
	}
}

func (forward *GostClientForward) GenerateRLimiter() config.LimiterConfig {
	return config.LimiterConfig{
		Name: "rLimiter_" + forward.Code,
		Limits: []string{
			fmt.Sprintf("$ %d", forward.RLimiter),
		},
	}
}

func (forward *GostClientForward) GenerateObs(host, nodeVersion string) config.ObserverConfig {
	if nodeVersion > "v1.1.2" {
		return config.ObserverConfig{}
	}
	return config.ObserverConfig{
		Name: "obs_" + forward.Code,
		Plugin: &config.PluginConfig{
			Type: "http",
			Addr: fmt.Sprintf("%s/api/v1/public/gost/obs?tunnel=%s", host, forward.Code),
		},
	}
}

func (forward *GostClientForward) GenerateWhiteAdmission() config.AdmissionConfig {
	if forward.WhiteEnable == 2 {
		return config.AdmissionConfig{}
	}
	return config.AdmissionConfig{
		Name:      "admission_white_" + forward.Code,
		Whitelist: true,
		Matchers:  forward.GetWhiteList(),
	}
}

func (forward *GostClientForward) GenerateBlackAdmission() config.AdmissionConfig {
	if forward.BlackEnable == 2 {
		return config.AdmissionConfig{}
	}
	return config.AdmissionConfig{
		Name:      "admission_black_" + forward.Code,
		Whitelist: false,
		Matchers:  forward.GetBlackList(),
	}
}
