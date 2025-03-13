package service

import (
	cache2 "github.com/patrickmn/go-cache"
	"server/service/common/cache"
)

type NodePortReq struct {
	Code string `json:"code"` // 节点编号
	Use  bool   `json:"use"`  // 1=被占用
	Port string `json:"port"` // 端口
}

func (service *service) NodePort(req NodePortReq) {
	cache.SetNodePortUse(req.Code, req.Port, req.Use, cache2.NoExpiration)
}
