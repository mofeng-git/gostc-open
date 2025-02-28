package service

import (
	"errors"
	"go.uber.org/zap"
	"server/model"
	"server/pkg/jwt"
	"server/repository"
	"server/service/gost_engine"
)

type EnableReq struct {
	Code   string `binding:"required" json:"code"`
	Enable int    `binding:"required" json:"enable"`
}

func (service *service) Enable(claims jwt.Claims, req EnableReq) error {
	db, _, log := repository.Get("")
	var user model.SystemUser
	if db.Where("code = ?", claims.Code).First(&user).RowsAffected == 0 {
		return errors.New("用户错误")
	}
	var tunnel model.GostClientTunnel
	if db.Preload("Node").Where("code = ? AND user_code = ?", req.Code, user.Code).First(&tunnel).RowsAffected == 0 {
		return errors.New("操作失败")
	}
	if tunnel.Enable == req.Enable {
		return nil
	}

	tunnel.Enable = req.Enable
	if err := db.Save(&tunnel).Error; err != nil {
		log.Error("启用或停用端口转发失败", zap.Error(err))
		return errors.New("操作失败")
	}
	if tunnel.Enable == 1 {
		gost_engine.ClientTunnelConfig(db, tunnel.Code)
	} else {
		gost_engine.ClientRemoveTunnelConfig(db, tunnel, tunnel.Node)
	}
	return nil
}
