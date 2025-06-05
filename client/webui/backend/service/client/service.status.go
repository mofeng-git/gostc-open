package service

import (
	"encoding/json"
	"errors"
	"gostc-sub/internal/common"
	service2 "gostc-sub/internal/service"
	rpcService "gostc-sub/internal/service/rpc"
	"gostc-sub/webui/backend/global"
	"gostc-sub/webui/backend/model"
)

type StatusReq struct {
	Key    string `binding:"required" json:"key"`
	Status int    `binding:"required" json:"status"`
}

func (*service) Status(req StatusReq) error {
	value, ok := global.ClientFS.Get(req.Key)
	if !ok {
		return errors.New("客户端不存在")
	}
	var client model.Client
	marshal, _ := json.Marshal(value)
	_ = json.Unmarshal(marshal, &client)

	svc, ok := global.ClientMap.Load(req.Key)
	if !ok {
		svc = rpcService.NewClient(common.GenerateWsUrl(client.Tls == 1, client.Address), req.Key)
		global.ClientMap.Store(req.Key, svc)
	}
	if req.Status == 1 {
		if err := svc.(*service2.Client).Start(); err != nil {
			return err
		}
	} else {
		svc.(*service2.Client).Stop()
	}
	return nil
}
