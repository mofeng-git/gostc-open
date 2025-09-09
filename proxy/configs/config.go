package configs

import (
	"encoding/json"
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"net"
	"os"
	"strings"
)

type Config struct {
	HTTPAddr    string `yaml:"http-addr"`
	HTTPSAddr   string `yaml:"https-addr"`
	Certs       string `yaml:"certs"`
	ApiAddr     string `yaml:"api-addr"`
	AutoSetCert bool   `yaml:"auto-set-cert"` // 自动申请证书

	Default DomainConfig            `yaml:"default"`
	Domains map[string]DomainConfig `yaml:"domains"`
}

type DomainConfig struct {
	Target       string `yaml:"target"`
	Cert         string `yaml:"cert"`
	Key          string `yaml:"key"`
	ForceHttps   bool   `yaml:"force-https"`
	DisableHttps bool   `yaml:"disable-https"`
}

var adapter = caddyconfig.GetAdapter("caddyfile")

func (c *Config) GenerateCaddyfile() string {
	_, httpPort, _ := net.SplitHostPort(c.HTTPAddr)
	_, httpsPort, _ := net.SplitHostPort(c.HTTPSAddr)
	var cfgs []string

	cfgs = append(cfgs, `
{ 
    admin off
}
`)

	// 处理默认页面
	if httpPort != "" {
		if c.Default.ForceHttps && httpsPort != "" && !c.Default.DisableHttps {
			cfgs = append(cfgs, generateHttpToHttpsCfg("", httpPort, fmt.Sprintf("https://:%s", httpsPort)))
		} else {
			cfgs = append(cfgs, generateHttpCfg("", httpPort, c.Default.Target))
		}
	}
	// 默认反代目标，不自动申请SSL
	if httpsPort != "" && !c.Default.DisableHttps {
		cfgs = append(cfgs, generateHttpsCfg("", httpsPort, c.Default.Target, c.Default.Cert, c.Default.Key, false))
	}

	// 处理其他站点
	for domain, target := range c.Domains {
		if strings.HasPrefix(domain, ".") {
			domain = "*" + domain
		}

		if httpPort != "" {
			if target.ForceHttps && httpsPort != "" && !target.DisableHttps {
				cfgs = append(cfgs, generateHttpToHttpsCfg(domain, httpPort, fmt.Sprintf("https://%s:%s", domain, httpsPort)))
			} else {
				cfgs = append(cfgs, generateHttpCfg(domain, httpPort, target.Target))
			}
		}
		if httpsPort != "" && !target.DisableHttps {
			cfgs = append(cfgs, generateHttpsCfg(domain, httpsPort, target.Target, target.Cert, target.Key, c.AutoSetCert))
		}
	}
	return strings.Join(cfgs, "\n")
}

func (c *Config) ParseCaddyFileConfig() ([]byte, []caddyconfig.Warning, error) {
	caddyfile := c.GenerateCaddyfile()
	result, warnMsg, err := adapter.Adapt([]byte(caddyfile), nil)
	return result, warnMsg, err
}

var caddyFileHttpTemplate = `
http://%s:%s {
	encode gzip
	reverse_proxy %s {
        flush_interval 0s
		header_up X-Real-IP {remote}
    }
}
`

func generateHttpCfg(domain, port, target string) string {
	return fmt.Sprintf(caddyFileHttpTemplate, domain, port, target)
}

var caddyFileHttpToHttpsTemplate = `
http://%s:%s {
	redir %s{uri} 302
}
`

func generateHttpToHttpsCfg(domain, port, target string) string {
	return fmt.Sprintf(caddyFileHttpToHttpsTemplate, domain, port, target)
}

var caddyFileHttpsTemplate = `
https://%s:%s {
	encode gzip
	tls %s %s
	reverse_proxy %s {
        flush_interval 0s
		header_up X-Real-IP {remote}
    }
}
`

var caddyFileHttpsCertTemplate = `
https://%s:%s {
	encode gzip
	reverse_proxy %s {
        flush_interval 0s
		header_up X-Real-IP {remote}
    }
}
`

func generateHttpsCfg(domain, port, target, certPath, keyPath string, autoSetCert bool) string {
	if certPath == "" || keyPath == "" {
		certPath = "internal"
		keyPath = ""
	} else {
		if stat, err := os.Stat(certPath); err != nil || stat.IsDir() {
			certPath = "internal"
			keyPath = ""
		}
		if stat, err := os.Stat(keyPath); err != nil || stat.IsDir() {
			certPath = "internal"
			keyPath = ""
		}
	}
	// 启用了自动申请证书，并且未手动设置证书，则启用自动申请证书
	// 泛域名，不自动申请SSL
	if autoSetCert && certPath == "internal" && !strings.HasPrefix(domain, "*") {
		return fmt.Sprintf(caddyFileHttpsCertTemplate, domain, port, target)
	} else {
		return fmt.Sprintf(caddyFileHttpsTemplate, domain, port, certPath, keyPath, target)
	}
}

func GenerateCaddyConfig(cfg []byte) (*caddy.Config, error) {
	var rootCfg map[string]json.RawMessage
	if err := json.Unmarshal(cfg, &rootCfg); err != nil {
		return nil, err
	}
	config := caddy.Config{
		Admin: &caddy.AdminConfig{
			Disabled: true,
		},
		Logging:    nil,
		StorageRaw: nil,
		AppsRaw:    make(map[string]json.RawMessage),
	}
	if apps, ok := rootCfg["apps"]; ok {
		if err := json.Unmarshal(apps, &config.AppsRaw); err != nil {
			return nil, err
		}
	}
	return &config, nil
}
