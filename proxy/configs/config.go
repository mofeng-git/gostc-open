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
	HTTPAddr  string `yaml:"http-addr"`
	HTTPSAddr string `yaml:"https-addr"`
	Certs     string `yaml:"certs"`
	ApiAddr   string `yaml:"api-addr"`

	Default DomainConfig            `yaml:"default"`
	Domains map[string]DomainConfig `yaml:"domains"`
}

type DomainConfig struct {
	Target     string `yaml:"target"`
	Cert       string `yaml:"cert"`
	Key        string `yaml:"key"`
	ForceHttps bool   `yaml:"force-https"`
}

var adapter = caddyconfig.GetAdapter("caddyfile")

func (c *Config) ParseCaddyFileConfig() ([]byte, []caddyconfig.Warning, error) {
	_, httpPort, _ := net.SplitHostPort(c.HTTPAddr)
	_, httpsPort, _ := net.SplitHostPort(c.HTTPSAddr)
	var cfgs []string

	cfgs = append(cfgs, `
{ 
	auto_https off
}
`)

	// 处理默认页面
	if httpPort != "" {
		if c.Default.ForceHttps && httpsPort != "" {
			cfgs = append(cfgs, generateHttpToHttpsCfg("", httpPort, fmt.Sprintf("https://:%s", httpsPort)))
		} else {
			cfgs = append(cfgs, generateHttpCfg("", httpPort, c.Default.Target))
		}
	}
	if httpsPort != "" {
		cfgs = append(cfgs, generateHttpsCfg("", httpsPort, c.Default.Target, c.Default.Cert, c.Default.Key))
	}

	// 处理其他站点
	for domain, target := range c.Domains {
		if strings.HasPrefix(domain, ".") {
			domain = "*" + domain
		}

		if httpPort != "" {
			if target.ForceHttps && httpsPort != "" {
				cfgs = append(cfgs, generateHttpToHttpsCfg(domain, httpPort, fmt.Sprintf("https://%s:%s", domain, httpsPort)))
			} else {
				cfgs = append(cfgs, generateHttpCfg(domain, httpPort, target.Target))
			}
		}
		if httpsPort != "" {
			cfgs = append(cfgs, generateHttpsCfg(domain, httpsPort, target.Target, target.Cert, target.Key))
		}
	}
	result, warnMsg, err := adapter.Adapt([]byte(strings.Join(cfgs, "\n")), nil)
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

func generateHttpsCfg(domain, port, target, certPath, keyPath string) string {
	if certPath == "" || keyPath == "" {
		certPath = "internal"
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
	return fmt.Sprintf(caddyFileHttpsTemplate, domain, port, certPath, keyPath, target)
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
