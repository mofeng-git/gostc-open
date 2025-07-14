package common

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

func GenerateWsUrl(tls bool, address string) string {
	var scheme string
	if tls {
		scheme = "wss"
	} else {
		scheme = "ws"
	}
	return scheme + "://" + address
}

func GenerateHttpUrl(tls bool, address string) string {
	var scheme string
	if tls {
		scheme = "https"
	} else {
		scheme = "http"
	}
	return scheme + "://" + address
}

type GenerateUrl interface {
	WsUrl() string
	HttpUrl() string
}

type generate struct {
	tls     bool
	address string
	apiUrl  string
	wsUrl   string
	mu      *sync.Mutex
	expAt   int64
}

func (g *generate) WsUrl() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	if time.Now().UnixNano() > g.expAt {
		g.refresh()
	}
	return g.wsUrl
}

func (g *generate) HttpUrl() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	if time.Now().UnixNano() > g.expAt {
		g.refresh()
	}
	return g.apiUrl
}

func NewGenerateUrl(tls bool, address string) GenerateUrl {
	g := &generate{
		tls:     tls,
		address: address,
		apiUrl:  "",
		wsUrl:   "",
		mu:      &sync.Mutex{},
		expAt:   time.Now().UnixNano(),
	}
	g.refresh()
	return g
}

func (g *generate) refresh() {
	var scheme string
	if g.tls {
		scheme = "https"
	} else {
		scheme = "http"
	}
	result := parseUrl(scheme + "://" + g.address)
	parse, _ := url.Parse(result)
	{
		var scheme string
		if parse.Scheme == "http" {
			scheme = "http"
		} else {
			scheme = "https"
		}
		g.apiUrl = scheme + "://" + parse.Host
	}
	{
		var scheme string
		if parse.Scheme == "http" {
			scheme = "ws"
		} else {
			scheme = "wss"
		}
		g.wsUrl = scheme + "://" + parse.Host
	}
	g.expAt = time.Now().Add(time.Second * 30).UnixNano()
}

func parseUrl(reqUrl string) string {
	// 最多跟随10次重定向
	maxRedirects := 10
	// 禁用自动跳转
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for i := 0; i < maxRedirects; i++ {
		resp, err := client.Get(reqUrl)
		if err != nil {
			return reqUrl
		}
		resp.Body.Close()

		//fmt.Printf("跳转 %d：%s (状态码: %d)\n", i+1, reqUrl, resp.StatusCode)

		// 检查是否为重定向
		if resp.StatusCode == http.StatusMovedPermanently || // 301
			resp.StatusCode == http.StatusFound || // 302
			resp.StatusCode == http.StatusSeeOther || // 303
			resp.StatusCode == http.StatusTemporaryRedirect || // 307
			resp.StatusCode == http.StatusPermanentRedirect { // 308

			location := resp.Header.Get("Location")
			if location == "" {
				return reqUrl
			}

			// 解析新的URL
			newURL, err := url.Parse(location)
			if err != nil {
				return reqUrl
			}
			// 处理相对URL
			baseURL, _ := url.Parse(reqUrl)
			reqUrl = baseURL.ResolveReference(newURL).String()
		} else {
			// 非重定向，结束跳转链
			return reqUrl
		}
	}
	return reqUrl
}
