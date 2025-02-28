package service

import (
	"net/url"
	"server/service/common/cache"
)

type IngressResp struct {
	Endpoint string `json:"endpoint"`
}

func (service *service) Ingress(host string) IngressResp {
	unescape, _ := url.QueryUnescape(host)
	return IngressResp{Endpoint: cache.GetIngress(unescape)}
}
