package service

import (
	"encoding/json"
	"errors"
	"gostc-sub/internal/common"
	service2 "gostc-sub/internal/service"
	service3 "gostc-sub/internal/service/visitor"
	"gostc-sub/webui/backend/global"
	"gostc-sub/webui/backend/model"
	"strconv"
)

type StatusReq struct {
	Key    string `binding:"required" json:"key"`
	Status int    `binding:"required" json:"status"`
}

func (*service) Status(req StatusReq) error {
	value, ok := global.TunnelFS.Get(req.Key)
	if !ok {
		return errors.New("私有隧道不存在")
	}
	var tunnel model.Tunnel
	marshal, _ := json.Marshal(value)
	_ = json.Unmarshal(marshal, &tunnel)

	svc, ok := global.TunnelMap.Load(req.Key)
	if !ok {
		port, _ := strconv.Atoi(tunnel.Port)
		generate := common.NewGenerateUrl(tunnel.Tls == 1, tunnel.Address)
		svc = service3.NewTunnel(generate, tunnel.Key, tunnel.Bind, port)
		global.TunnelMap.Store(req.Key, svc)
	}
	if req.Status == 1 {
		if err := svc.(service2.Service).Start(); err != nil {
			return err
		}
	} else {
		svc.(service2.Service).Stop()
	}
	return nil
}
