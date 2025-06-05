package proxy

import (
	"crypto/tls"
	"net/http/httputil"
	"net/url"
)

type Config struct {
	HTTPAddr  string                  `yaml:"http-addr"`
	HTTPSAddr string                  `yaml:"https-addr"`
	Default   DomainConfig            `yaml:"default"`
	Domains   map[string]DomainConfig `yaml:"domains"`
}

type DomainConfig struct {
	Target     string `yaml:"target"`
	Cert       string `yaml:"cert"`
	Key        string `yaml:"key"`
	ForceHttps bool   `yaml:"force-https"` // 是否强制HTTPS
}

func (cfg DomainConfig) Generate() (*url.URL, *httputil.ReverseProxy, *tls.Certificate) {
	target, _ := url.Parse(cfg.Target)
	if target == nil {
		return nil, nil, nil
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	if cfg.Cert != "" && cfg.Key != "" {
		cert, err := tls.LoadX509KeyPair(cfg.Cert, cfg.Key)
		if err != nil {
			return target, proxy, nil
		}
		return target, proxy, &cert
	}
	return target, proxy, nil
}
