package service

import (
	service2 "gostc-sub/internal/service"
	"gostc-sub/webui/backend/global"
)

type DeleteReq struct {
	Key string `binding:"required" json:"key"`
}

func (*service) Delete(req DeleteReq) error {
	if err := global.TunnelFS.Delete(req.Key); err != nil {
		return err
	}
	tunnel, ok := global.TunnelMap.Load(req.Key)
	if ok {
		tunnel.(*service2.Tunnel).Stop()
	}
	return nil
}
