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

	proxy, _ := db.GostClientProxy.Preload(db.GostClientProxy.Node).Where(
		db.GostClientProxy.Code.Eq(req.Code),
		db.GostClientProxy.UserCode.Eq(user.Code),
	).First()
	if proxy == nil {
		return errors.New("操作失败")
	}
	if proxy.Enable == req.Enable {
		return nil
	}

	proxy.Enable = req.Enable
	if err := db.GostClientProxy.Save(proxy); err != nil {
		log.Error("启用或停用代理隧道失败", zap.Error(err))
		return errors.New("操作失败")
	}
	if proxy.Enable == 1 {
		gost_engine.ClientProxyConfig(db, proxy.Code)
	} else {
		gost_engine.ClientRemoveProxyConfig(*proxy, proxy.Node)
	}
	return nil
}
