package proxy

import (
	"crypto/tls"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Server struct {
	mu            sync.RWMutex
	domainMap     map[string]*httputil.ReverseProxy
	targetMap     map[string]*url.URL
	certMap       map[string]*tls.Certificate
	defaultCert   *tls.Certificate
	defaultProxy  *httputil.ReverseProxy
	defaultTarget *url.URL
	httpAddr      string
	httpsAddr     string
	httpServer    *http.Server
	httpsServer   *http.Server
	log           *zap.Logger
}

func NewServer(cfg Config, log *zap.Logger) *Server {
	server := Server{
		mu:            sync.RWMutex{},
		domainMap:     make(map[string]*httputil.ReverseProxy),
		targetMap:     make(map[string]*url.URL),
		certMap:       make(map[string]*tls.Certificate),
		defaultCert:   nil,
		defaultProxy:  nil,
		defaultTarget: nil,
		httpAddr:      cfg.HTTPAddr,
		httpsAddr:     cfg.HTTPSAddr,
		log:           log,
	}

	// 加载域名配置
	for domain, domainConfig := range cfg.Domains {
		target, err := url.Parse(domainConfig.Target)
		if err != nil {
			server.log.Warn(fmt.Sprintf("parse target fail,%s", domainConfig.Target), zap.Error(err))
			continue
		}
		proxy := httputil.NewSingleHostReverseProxy(target)
		configureTransport(proxy)
		server.domainMap[domain] = proxy
		server.targetMap[domain] = target
		if domainConfig.Cert != "" && domainConfig.Key != "" {
			cert, err := tls.LoadX509KeyPair(domainConfig.Cert, domainConfig.Key)
			if err == nil {
				server.certMap[domain] = &cert
			} else {
				server.log.Warn(fmt.Sprintf("Load %s,%s cert failed", domainConfig.Cert, domainConfig.Key), zap.Error(err))
			}
		}
	}

	// 加载默认配置
	target, err := url.Parse(cfg.Default.Target)
	if err != nil {
		server.log.Warn(fmt.Sprintf("parse default target fail,%s", cfg.Default.Target), zap.Error(err))
	} else {
		server.defaultTarget = target
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	configureTransport(proxy)
	server.defaultProxy = proxy
	if cfg.Default.Cert != "" && cfg.Default.Key != "" {
		cert, err := tls.LoadX509KeyPair(cfg.Default.Cert, cfg.Default.Key)
		if err == nil {
			server.defaultCert = &cert
		} else {
			server.log.Warn(fmt.Sprintf("Load %s,%s default cert failed", cfg.Default.Cert, cfg.Default.Key), zap.Error(err))
		}
	}
	if server.defaultCert == nil {
		cert, err := genCertificate(time.Hour*24*3650, "sian", "gost.sian.one")
		if err != nil {
			server.log.Warn(fmt.Sprintf("gen default cert failed"), zap.Error(err))
		} else {
			server.defaultCert = &cert
		}
	}
	return &server
}

func (server *Server) UpdateDomain(domain, target, certFile, keyFile string) {
	server.mu.Lock()
	defer server.mu.Unlock()
	if domain == "" {
		return
	}

	targetUrl, err := url.Parse(target)
	if err != nil {
		server.log.Warn(fmt.Sprintf("parse target fail,%s", targetUrl), zap.Error(err))
		return
	}
	server.targetMap[domain] = targetUrl

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	configureTransport(proxy)
	server.domainMap[domain] = proxy

	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			server.log.Warn(fmt.Sprintf("Load %s,%s default cert failed", certFile, keyFile), zap.Error(err))
		} else {
			server.certMap[domain] = &cert
		}
	}
}

func configureTransport(proxy *httputil.ReverseProxy) {
	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// 自定义Director保留原始Host头
	originalDirector := proxy.Director
	proxy.Director = func(r *http.Request) {
		originalDirector(r)
		// 从Header恢复原始Host
		if forwardedHost := r.Header.Get("X-Forwarded-Host"); forwardedHost != "" {
			r.Host = forwardedHost
		}
		// 确保Upgrade头传递
		if r.Header.Get("Upgrade") == "websocket" {
			r.Header.Set("Connection", "Upgrade")
		}
	}

	// 处理重定向响应
	proxy.ModifyResponse = func(resp *http.Response) error {
		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			handleRedirectLocation(resp)
		}
		return nil
	}
}

func handleRedirectLocation(resp *http.Response) {
	location, err := resp.Location()
	if err != nil || location == nil {
		return
	}

	// 获取原始请求的协议和Host
	req := resp.Request
	scheme := "http"
	if req.TLS != nil {
		scheme = "https"
	}
	originalHost := req.Host

	// 修正Location为代理地址
	if location.Host != originalHost {
		location.Scheme = scheme
		location.Host = originalHost
		resp.Header.Set("Location", location.String())
	}
}

func (server *Server) StartHTTPServer() error {
	server.httpServer = &http.Server{
		Addr:    server.httpAddr,
		Handler: http.HandlerFunc(server.httpProxyHandler),
	}
	return server.httpServer.ListenAndServe()
}

func (server *Server) StartHTTPSServer() error {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		GetCertificate: func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			server.mu.RLock()
			defer server.mu.RUnlock()
			// 查找匹配的证书
			var domain = hello.ServerName
			cert, exists := server.certMap[domain]
			if exists {
				return cert, nil
			}

			domain = "." + strings.Join(strings.Split(domain, ".")[1:], ".")
			cert, exists = server.certMap[domain]
			if exists {
				return cert, nil
			}
			domain = "*." + strings.Join(strings.Split(domain, ".")[1:], ".")
			cert, exists = server.certMap[domain]
			if exists {
				return cert, nil
			}

			if server.defaultCert == nil {
				server.log.Warn("no match cert", zap.String("serverName", hello.ServerName))
				return nil, errors.New("no match cert")
			}
			return server.defaultCert, nil
		},
	}

	server.httpsServer = &http.Server{
		Addr:      server.httpsAddr,
		Handler:   http.HandlerFunc(server.httpsProxyHandler),
		TLSConfig: tlsConfig,
	}

	return server.httpsServer.ListenAndServeTLS("", "") // 证书通过GetCertificate动态加载
}

func (server *Server) findProxy(domain string) (proxy *httputil.ReverseProxy, err error) {
	server.mu.RLock()
	defer server.mu.RUnlock()
	var exists bool
	proxy, exists = server.domainMap[domain]
	if exists {
		return proxy, nil
	}

	domain = "." + strings.Join(strings.Split(domain, ".")[1:], ".")
	proxy, exists = server.domainMap[domain]
	if exists {
		return proxy, nil
	}
	domain = "*." + strings.Join(strings.Split(domain, ".")[1:], ".")
	proxy, exists = server.domainMap[domain]
	if exists {
		return proxy, nil
	}

	if server.defaultProxy == nil {
		return nil, errors.New("no match backend for this domain")
	}

	return server.defaultProxy, nil
}
func (server *Server) findTarget(domain string) (target *url.URL, err error) {
	server.mu.RLock()
	defer server.mu.RUnlock()
	var exists bool
	target, exists = server.targetMap[domain]
	if exists {
		return target, nil
	}

	domain = "." + strings.Join(strings.Split(domain, ".")[1:], ".")
	target, exists = server.targetMap[domain]
	if exists {
		return target, nil
	}
	domain = "*." + strings.Join(strings.Split(domain, ".")[1:], ".")
	target, exists = server.targetMap[domain]
	if exists {
		return target, nil
	}

	if server.defaultTarget == nil {
		return nil, errors.New("no match backend for this domain")
	}

	return server.defaultTarget, nil
}

func (server *Server) httpProxyHandler(w http.ResponseWriter, r *http.Request) {
	// 设置代理相关Header
	r.Header.Set("X-Forwarded-Host", r.Host)
	r.Header.Set("X-Forwarded-Proto", "http") // 标记原始协议

	// 保留WebSocket头
	if strings.EqualFold(r.Header.Get("Connection"), "upgrade") {
		r.Header.Set("Connection", "Upgrade")
	}

	host := r.Host
	if hostHeader := r.Header.Get("Host"); hostHeader != "" {
		host = hostHeader
	}
	domain := host
	if h, _, err := net.SplitHostPort(host); err == nil {
		domain = h
	}

	proxy, err := server.findProxy(domain)
	if err != nil {
		http.Error(w, "No backend configured for this domain", http.StatusBadGateway)
		return
	}

	target, err := server.findTarget(domain)
	if err != nil {
		http.Error(w, "No backend configured for this domain", http.StatusBadGateway)
		return
	}

	r.Header.Set("X-Forwarded-Host", r.Host)
	r.Header.Set("X-Origin-Host", target.Host)
	r.Host = domain
	proxy.ServeHTTP(w, r)
}

func (server *Server) httpsProxyHandler(w http.ResponseWriter, r *http.Request) {
	// 设置代理相关Header
	r.Header.Set("X-Forwarded-Host", r.Host)
	r.Header.Set("X-Forwarded-Proto", "https") // 标记原始协议

	// 保留WebSocket头
	if strings.EqualFold(r.Header.Get("Connection"), "upgrade") {
		r.Header.Set("Connection", "Upgrade")
	}

	host := r.Host
	if hostHeader := r.Header.Get("Host"); hostHeader != "" {
		host = hostHeader
	}
	domain := host
	if h, _, err := net.SplitHostPort(host); err == nil {
		domain = h
	}

	proxy, err := server.findProxy(domain)
	if err != nil {
		http.Error(w, "No backend configured for this domain", http.StatusBadGateway)
		return
	}

	target, err := server.findTarget(domain)
	if err != nil {
		http.Error(w, "No backend configured for this domain", http.StatusBadGateway)
		return
	}

	r.Header.Set("X-Forwarded-Host", r.Host)
	r.Header.Set("X-Origin-Host", target.Host)
	r.Host = domain
	proxy.ServeHTTP(w, r)
}

func (server *Server) Close() {
	if server.httpServer != nil {
		_ = server.httpServer.Close()
	}
	if server.httpsServer != nil {
		_ = server.httpsServer.Close()
	}
}
