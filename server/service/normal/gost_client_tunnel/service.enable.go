package service

import (
	"errors"
	"go.uber.org/zap"
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
	user, _ := db.SystemUser.Where(db.SystemUser.Code.Eq(claims.Code)).First()
	if user == nil {
		return errors.New("用户错误")
	}

	tunnel, _ := db.GostClientTunnel.Where(db.GostClientTunnel.Code.Eq(req.Code), db.GostClientTunnel.UserCode.Eq(claims.Code)).First()
	if tunnel == nil {
		return errors.New("操作失败")
	}
	if tunnel.Enable == req.Enable {
		return nil
	}

	tunnel.Enable = req.Enable
	if err := db.GostClientTunnel.Save(tunnel); err != nil {
		log.Error("启用或停用端口转发失败", zap.Error(err))
		return errors.New("操作失败")
	}
	if tunnel.Enable == 1 {
		gost_engine.ClientTunnelConfig(db, tunnel.Code)
	} else {
		gost_engine.ClientRemoveTunnelConfig(db, *tunnel, tunnel.Node)
	}
	return nil
}
