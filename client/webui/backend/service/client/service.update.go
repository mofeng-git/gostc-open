package service

import (
	"errors"
	"gostc-sub/internal/common"
	service2 "gostc-sub/internal/service"
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
	if common.State.Get(req.Key) {
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
	client := service2.NewClient(common.GenerateWsUrl(req.Tls == 1, req.Address), req.Key)
	global.ClientMap.Store(req.Key, client)
	return nil
}
