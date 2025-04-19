package service

import (
	"errors"
	"gostc-sub/internal/common"
	service2 "gostc-sub/internal/service"
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
	if utils.IsUse(req.Bind, port) {
		return errors.New("本地端口已被占用")
	}
	if err := global.TunnelFS.Insert(req.Key, model.Tunnel{
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
	tunnel := service2.NewTunnel(common.GenerateHttpUrl(req.Tls == 1, req.Address), req.Key, req.Bind, req.Port)
	global.TunnelMap.Store(req.Key, tunnel)
	return nil
}
