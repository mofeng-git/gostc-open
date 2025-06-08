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

	host, _ := db.GostClientHost.Where(
		db.GostClientHost.UserCode.Eq(user.Code),
		db.GostClientHost.Code.Eq(req.Code),
	).First()
	if host == nil {
		return errors.New("操作失败")
	}
	if host.Enable == req.Enable {
		return nil
	}

	host.Enable = req.Enable
	if err := db.GostClientHost.Save(host); err != nil {
		log.Error("启用或停用域名解析失败", zap.Error(err))
		return errors.New("操作失败")
	}
	if host.Enable == 1 {
		engine.ClientHostConfig(db, host.Code)
	} else {
		engine.ClientRemoveHostConfig(db, *host, host.Node)
	}
	return nil
}
