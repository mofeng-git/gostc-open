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
	mu                sync.RWMutex
	domainMatchersMap map[string]*domainMatcher
	proxyMap          map[string]*httputil.ReverseProxy
	targetMap         map[string]*url.URL
	certMap           map[string]*tls.Certificate
	forceHttps        map[string]*bool
	defaultCert       *tls.Certificate
	defaultProxy      *httputil.ReverseProxy
	defaultTarget     *url.URL
	httpAddr          string
	httpsAddr         string
	httpServer        *http.Server
	httpsServer       *http.Server
	log               *zap.Logger
}

type domainMatcher struct {
	FullDomain string // example: www.example.com
	DotDomain  string // example: .example.com
	StarDomain string // example: *.example.com
}

func NewServer(cfg Config, log *zap.Logger) *Server {
	server := Server{
		mu:                sync.RWMutex{},
		domainMatchersMap: make(map[string]*domainMatcher),
		proxyMap:          make(map[string]*httputil.ReverseProxy),
		targetMap:         make(map[string]*url.URL),
		certMap:           make(map[string]*tls.Certificate),
		forceHttps:        make(map[string]*bool),
		defaultCert:       nil,
		defaultProxy:      nil,
		defaultTarget:     nil,
		httpAddr:          cfg.HTTPAddr,
		httpsAddr:         cfg.HTTPSAddr,
		httpServer:        nil,
		httpsServer:       nil,
		log:               log,
	}

	// 加载域名配置
	for domain, domainConfig := range cfg.Domains {
		target, proxy, cert := domainConfig.Generate()
		if target != nil {
			server.targetMap[domain] = target
			if proxy != nil {
				server.configureTransport(proxy)
				server.proxyMap[domain] = proxy
			}
			if cert != nil {
				server.certMap[domain] = cert
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

	target, proxy, cert := cfg.Default.Generate()
	if target != nil {
		server.defaultTarget = target
		if proxy != nil {
			server.configureTransport(proxy)
			server.defaultProxy = proxy
		}
		if cert != nil {
			server.defaultCert = cert
		} else {
			defaultCert, _ := genCertificate(time.Hour*24*3650, "sian", "gost.sian.one")
			server.defaultCert = &defaultCert
		}
	}
	return &server
}

func (server *Server) UpdateDomain(domain string, cfg DomainConfig) {
	server.mu.Lock()
	defer server.mu.Unlock()
	if domain == "" {
		return
	}
	target, proxy, cert := cfg.Generate()
	if target != nil {
		server.targetMap[domain] = target
		if proxy != nil {
			server.configureTransport(proxy)
			server.proxyMap[domain] = proxy
		}
		if cert != nil {
			server.certMap[domain] = cert
		} else {
			delete(server.certMap, domain)
		}
	}
	server.forceHttps[domain] = &cfg.ForceHttps
}

func (server *Server) configureTransport(proxy *httputil.ReverseProxy) {
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
		if r.TLS != nil {
			r.Header.Set("X-Forwarded-Proto", "https")
		} else {
			r.Header.Set("X-Forwarded-Proto", "http")
		}
		// 确保Upgrade头传递
		if r.Header.Get("Upgrade") == "websocket" {
			r.Header.Set("Connection", "Upgrade")
		}
	}

	// 处理重定向响应
	proxy.ModifyResponse = func(resp *http.Response) error {
		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			server.handleRedirectLocation(resp)
		}
		return nil
	}
}

func (server *Server) handleRedirectLocation(resp *http.Response) {
	location, err := resp.Location()
	if err != nil || location == nil {
		return
	}
	req := resp.Request

	// 仅修改与内网地址一致的重定向location
	reqSplit := strings.Split(req.Host, ":")
	locationSplit := strings.Split(location.Host, ":")
	if reqSplit[0] == locationSplit[0] {
		// 获取原始请求的协议
		scheme := "http"
		if req.TLS != nil || req.Header.Get("X-Forwarded-Proto") == "https" {
			scheme = "https"
		}
		// 创建新的URL，保留路径和查询参数
		newLocation := &url.URL{
			Scheme:   scheme,
			Host:     req.Host,
			Path:     location.Path,
			RawQuery: location.RawQuery,
			Fragment: location.Fragment,
		}
		resp.Header.Set("Location", newLocation.String())
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
			matcher := server.findDomainMatcher(hello.ServerName)
			cert, err := server.findCert(matcher)
			if err != nil {
				return nil, err
			}
			return cert, nil
		},
	}

	server.httpsServer = &http.Server{
		Addr:      server.httpsAddr,
		Handler:   http.HandlerFunc(server.httpsProxyHandler),
		TLSConfig: tlsConfig,
	}

	return server.httpsServer.ListenAndServeTLS("", "") // 证书通过GetCertificate动态加载
}

func (server *Server) findDomainMatcher(domain string) *domainMatcher {
	server.mu.RLock()
	matcher, ok := server.domainMatchersMap[domain]
	if ok {
		server.mu.RUnlock()
		return matcher
	}
	server.mu.RUnlock()

	server.mu.Lock()
	matcher = &domainMatcher{
		FullDomain: domain,
		DotDomain:  "." + strings.Join(strings.Split(domain, ".")[1:], "."),
		StarDomain: "*." + strings.Join(strings.Split(domain, ".")[1:], "."),
	}
	server.domainMatchersMap[domain] = matcher
	server.mu.Unlock()
	return matcher
}

func (server *Server) findCert(matcher *domainMatcher) (cert *tls.Certificate, err error) {
	server.mu.RLock()
	defer server.mu.RUnlock()
	var exists bool
	cert, exists = server.certMap[matcher.FullDomain]
	if exists {
		return cert, nil
	}
	cert, exists = server.certMap[matcher.DotDomain]
	if exists {
		return cert, nil
	}
	cert, exists = server.certMap[matcher.StarDomain]
	if exists {
		return cert, nil
	}
	if server.defaultCert == nil {
		return nil, errors.New("no match backend for this domain")
	}
	return server.defaultCert, nil
}

func (server *Server) findProxy(matcher *domainMatcher) (proxy *httputil.ReverseProxy, err error) {
	server.mu.RLock()
	defer server.mu.RUnlock()
	var exists bool
	proxy, exists = server.proxyMap[matcher.FullDomain]
	if exists {
		return proxy, nil
	}
	proxy, exists = server.proxyMap[matcher.DotDomain]
	if exists {
		return proxy, nil
	}
	proxy, exists = server.proxyMap[matcher.StarDomain]
	if exists {
		return proxy, nil
	}
	if server.defaultProxy == nil {
		return nil, errors.New("no match backend for this domain")
	}
	return server.defaultProxy, nil
}

func (server *Server) findTarget(matcher *domainMatcher) (target *url.URL, err error) {
	server.mu.RLock()
	defer server.mu.RUnlock()
	var exists bool
	target, exists = server.targetMap[matcher.FullDomain]
	if exists {
		return target, nil
	}
	target, exists = server.targetMap[matcher.DotDomain]
	if exists {
		return target, nil
	}
	target, exists = server.targetMap[matcher.StarDomain]
	if exists {
		return target, nil
	}
	if server.defaultTarget == nil {
		return nil, errors.New("no match backend for this domain")
	}
	return server.defaultTarget, nil
}

func (server *Server) isForceHttps(matcher *domainMatcher) bool {
	if server.httpsServer == nil {
		return false
	}
	server.mu.RLock()
	defer server.mu.RUnlock()
	var exists bool
	var target *bool
	target, exists = server.forceHttps[matcher.FullDomain]
	if exists {
		return *target
	}
	target, exists = server.forceHttps[matcher.DotDomain]
	if exists {
		return *target
	}
	target, exists = server.forceHttps[matcher.StarDomain]
	if exists {
		return *target
	}
	return false
}

func (server *Server) httpsRedirect(w http.ResponseWriter, r *http.Request) {
	host, _, _ := net.SplitHostPort(r.Host)
	_, httpsPort, _ := net.SplitHostPort(server.httpsAddr)
	http.Redirect(w, r, "https://"+host+":"+httpsPort+r.RequestURI, http.StatusMovedPermanently)
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
	matcher := server.findDomainMatcher(domain)

	// 强制HTTPS
	if server.isForceHttps(matcher) {
		server.httpsRedirect(w, r)
		return
	}

	proxy, err := server.findProxy(matcher)
	if err != nil {
		http.Error(w, "No backend configured for this domain", http.StatusBadGateway)
		return
	}

	target, err := server.findTarget(matcher)
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

	matcher := server.findDomainMatcher(domain)

	proxy, err := server.findProxy(matcher)
	if err != nil {
		http.Error(w, "No backend configured for this domain", http.StatusBadGateway)
		return
	}
	target, err := server.findTarget(matcher)
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
