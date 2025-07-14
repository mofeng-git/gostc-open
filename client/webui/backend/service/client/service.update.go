package service

import (
	"errors"
	"gostc-sub/internal/common"
	rpcService "gostc-sub/internal/service"
	"gostc-sub/webui/backend/global"
	"gostc-sub/webui/backend/model"
)

type UpdateReq struct {
	Name      string `binding:"required" json:"name"`
	Address   string `binding:"required" json:"address"`
	Tls       int    `binding:"required" json:"tls"`
	Key       string `binding:"required" json:"key"`
	AutoStart int    `json:"autoStart"`
}

func (*service) Update(req UpdateReq) error {
	if rpcService.State.Get(req.Key) {
		return errors.New("客户端正在运行中，请停止运行后修改")
	}
	if err := global.ClientFS.Update(req.Key, model.Client{
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
