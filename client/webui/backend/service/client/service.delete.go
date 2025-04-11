package service

import (
	service2 "gostc-sub/internal/service"
	"gostc-sub/webui/backend/global"
)

type DeleteReq struct {
	Key string `binding:"required" json:"key"`
}

func (*service) Delete(req DeleteReq) error {
	if err := global.ClientFS.Delete(req.Key); err != nil {
		return err
	}
	client, ok := global.ClientMap.Load(req.Key)
	if ok {
		client.(*service2.Client).Stop()
	}
	return nil
}
