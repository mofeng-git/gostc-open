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

	forward, _ := db.GostClientForward.Preload(db.GostClientForward.Node).Where(
		db.GostClientForward.Code.Eq(req.Code),
		db.GostClientForward.UserCode.Eq(user.Code),
	).First()
	if forward == nil {
		return errors.New("操作失败")
	}
	if forward.Enable == req.Enable {
		return nil
	}

	forward.Enable = req.Enable
	if err := db.GostClientForward.Save(forward); err != nil {
		log.Error("启用或停用端口转发失败", zap.Error(err))
		return errors.New("操作失败")
	}
	if forward.Enable == 1 {
		gost_engine.ClientForwardConfig(db, forward.Code)
	} else {
		gost_engine.ClientRemoveForwardConfig(*forward, forward.Node)
	}
	return nil
}
