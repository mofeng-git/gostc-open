package service

import (
	"server/service/common/cache"
	"time"
)

type ClientPortReq struct {
	Code string `json:"code"` // 节点编号
	Use  bool   `json:"use"`  // 1=被占用
	Port string `json:"port"` // 端口
}

func (service *service) ClientPort(req ClientPortReq) {
	cache.SetClientPortUse(req.Code, req.Port, req.Use, time.Second*3)
}
