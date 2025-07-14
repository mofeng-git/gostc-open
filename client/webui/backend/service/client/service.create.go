package service

import (
	"gostc-sub/internal/common"
	rpcService "gostc-sub/internal/service"
	"gostc-sub/webui/backend/global"
	"gostc-sub/webui/backend/model"
)

type CreateReq struct {
	Name      string `binding:"required" json:"name"`
	Address   string `binding:"required" json:"address"`
	Tls       int    `binding:"required" json:"tls"`
	Key       string `binding:"required" json:"key"`
	AutoStart int    `json:"autoStart"`
}

func (*service) Create(req CreateReq) error {
	if err := global.ClientFS.Insert(req.Key, model.Client{
		Key:       req.Key,
		Name:      req.Name,
		Address:   req.Address,
		Tls:       req.Tls,
		AutoStart: req.AutoStart,
	}); err != nil {
		return err
	}
	generate := common.NewGenerateUrl(req.Tls == 1, req.Address)
	client := rpcService.NewClient(generate, req.Key)
	global.ClientMap.Store(req.Key, client)
	return nil
}
