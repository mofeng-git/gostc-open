package service

import (
	service2 "gostc-sub/internal/service"
	"gostc-sub/webui/backend/global"
)

type DeleteReq struct {
	Key string `binding:"required" json:"key"`
}

func (*service) Delete(req DeleteReq) error {
	if err := global.P2PFS.Delete(req.Key); err != nil {
		return err
	}
	p2p, ok := global.P2PMap.Load(req.Key)
	if ok {
		p2p.(service2.Service).Stop()
	}
	return nil
}
