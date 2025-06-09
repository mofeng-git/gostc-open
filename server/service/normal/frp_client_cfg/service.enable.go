package service

import (
	"errors"
	"go.uber.org/zap"
	"server/pkg/jwt"
	"server/repository"
	"server/repository/query"
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

	return db.Transaction(func(tx *query.Query) error {
		cfg, _ := tx.FrpClientCfg.Where(
			tx.FrpClientCfg.UserCode.Eq(user.Code),
			tx.FrpClientCfg.Code.Eq(req.Code),
		).First()
		if cfg == nil {
			return errors.New("操作失败")
		}
		if cfg.Enable == req.Enable {
			return nil
		}
		cfg.Enable = req.Enable
		if err := tx.FrpClientCfg.Save(cfg); err != nil {
			log.Error("启用或停用配置隧道失败", zap.Error(err))
			return errors.New("操作失败")
		}
		if cfg.Enable == 1 {
			return engine.ClientCfgConfig(tx, cfg.Code)
		} else {
			engine.ClientCfgRemove(cfg.ClientCode, cfg.Code)
			return nil
		}
	})
}
