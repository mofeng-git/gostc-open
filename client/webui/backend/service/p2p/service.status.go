package service

import (
	"encoding/json"
	"errors"
	"gostc-sub/internal/common"
	service2 "gostc-sub/internal/service"
	"gostc-sub/webui/backend/global"
	"gostc-sub/webui/backend/model"
)

type StatusReq struct {
	Key    string `binding:"required" json:"key"`
	Status int    `binding:"required" json:"status"`
}

func (*service) Status(req StatusReq) error {
	value, ok := global.P2PFS.Get(req.Key)
	if !ok {
		return errors.New("P2P隧道不存在")
	}
	var p2p model.P2P
	marshal, _ := json.Marshal(value)
	_ = json.Unmarshal(marshal, &p2p)

	svc, ok := global.P2PMap.Load(req.Key)
	if !ok {
		svc = service2.NewP2P(common.GenerateHttpUrl(p2p.Tls == 1, p2p.Address), p2p.Key, p2p.Bind, p2p.Port)
		global.P2PMap.Store(req.Key, svc)
	}
	if req.Status == 1 {
		if err := svc.(*service2.P2P).Start(); err != nil {
			return err
		}
	} else {
		svc.(*service2.P2P).Stop()
	}
	return nil
}
