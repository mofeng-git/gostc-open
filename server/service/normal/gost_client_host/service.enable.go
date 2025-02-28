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
	var host model.GostClientHost
	if db.Preload("Node").Where("code = ? AND user_code = ?", req.Code, user.Code).First(&host).RowsAffected == 0 {
		return errors.New("操作失败")
	}
	if host.Enable == req.Enable {
		return nil
	}

	host.Enable = req.Enable
	if err := db.Save(&host).Error; err != nil {
		log.Error("启用或停用域名解析失败", zap.Error(err))
		return errors.New("操作失败")
	}
	if host.Enable == 1 {
		gost_engine.ClientHostConfig(db, host.Code)
	} else {
		gost_engine.ClientRemoveHostConfig(db, host, host.Node)
	}
	return nil
}
