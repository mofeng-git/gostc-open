package service

import (
	"errors"
	"gostc-sub/internal/common"
	service3 "gostc-sub/internal/service/visitor"
	"gostc-sub/pkg/utils"
	"gostc-sub/webui/backend/global"
	"gostc-sub/webui/backend/model"
	"strconv"
)

type CreateReq struct {
	Name      string `binding:"required" json:"name"`
	Key       string `binding:"required" json:"key"`
	Bind      string `json:"bind"`
	Port      string `binding:"required" json:"port"`
	Address   string `binding:"required" json:"address"`
	Tls       int    `binding:"required" json:"tls"`
	AutoStart int    `json:"autoStart"`
}

func (*service) Create(req CreateReq) error {
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

	if err := global.P2PFS.Insert(req.Key, model.P2P{
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
