package service

import (
	"errors"
	"go.uber.org/zap"
	"server/pkg/jwt"
	"server/repository"
	"server/service/engine"
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

	p2p, _ := db.GostClientP2P.Where(db.GostClientP2P.Code.Eq(req.Code), db.GostClientP2P.UserCode.Eq(claims.Code)).First()
	if p2p == nil {
		return errors.New("操作失败")
	}
	if p2p.Enable == req.Enable {
		return nil
	}

	p2p.Enable = req.Enable
	if err := db.GostClientP2P.Save(p2p); err != nil {
		log.Error("启用或停用端口转发失败", zap.Error(err))
		return errors.New("操作失败")
	}
	if p2p.Enable == 1 {
		engine.ClientP2PConfig(db, p2p.Code)
	} else {
		engine.ClientRemoveP2PConfig(*p2p)
	}
	return nil
}
