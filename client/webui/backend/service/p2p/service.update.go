package service

import (
	"errors"
	"gostc-sub/internal/common"
	service2 "gostc-sub/internal/service"
	service3 "gostc-sub/internal/service/visitor"
	"gostc-sub/pkg/utils"
	"gostc-sub/webui/backend/global"
	"gostc-sub/webui/backend/model"
	"strconv"
)

type UpdateReq struct {
	Name      string `binding:"required" json:"name"`
	Key       string `binding:"required" json:"key"`
	Bind      string `json:"bind"`
	Port      string `binding:"required" json:"port"`
	Tls       int    `binding:"required" json:"tls"`
	Address   string `binding:"required" json:"address"`
	AutoStart int    `json:"autoStart"`
}

func (*service) Update(req UpdateReq) error {
	if !utils.ValidatePort(req.Port) {
		return errors.New("本地端口格式错误")
	}
	port, err := strconv.Atoi(req.Port)
	if err != nil {
		return errors.New("端口格式错误")
	}
	if err := utils.IsUse(req.Bind, port); err != nil {
		return err
	}
	if service2.State.Get(req.Key) {
		return errors.New("P2P隧道正在运行中，请停止运行后修改")
	}
	if err := global.P2PFS.Update(req.Key, model.P2P{
		Key:       req.Key,
		Name:      req.Name,
		Bind:      req.Bind,
		Port:      req.Port,
		Address:   req.Address,
		Tls:       req.Tls,
		AutoStart: req.AutoStart,
	}); err != nil {
		return err
	}
	generate := common.NewGenerateUrl(req.Tls == 1, req.Address)
	p2p := service3.NewP2P(generate, req.Key, req.Bind, port)
	global.P2PMap.Store(req.Key, p2p)
	return nil
}
