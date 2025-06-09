package service

import (
	"errors"
	"go.uber.org/zap"
	"server/pkg/jwt"
	"server/repository"
	"server/repository/query"
	"server/service/engine"
)

type MigrateReq struct {
	Code       string `binding:"required" json:"code"`
	ClientCode string `binding:"required" json:"clientCode"`
}

func (service *service) Migrate(claims jwt.Claims, req MigrateReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(claims.Code)).First()
		if user == nil {
			return errors.New("用户错误")
		}
		newClient, _ := tx.GostClient.Where(tx.GostClient.UserCode.Eq(claims.Code), tx.GostClient.Code.Eq(req.ClientCode)).First()
		if newClient == nil {
			return errors.New("新客户端不存在")
		}
		cfg, _ := tx.FrpClientCfg.Where(
			tx.FrpClientCfg.UserCode.Eq(user.Code),
			tx.FrpClientCfg.Code.Eq(req.Code),
		).First()
		if cfg == nil {
			return errors.New("操作失败")
		}
		if cfg.ClientCode == req.ClientCode {
			return nil
		}
		if cfg.Enable == 1 {
			engine.ClientCfgRemove(cfg.ClientCode, cfg.Code)
		}
		cfg.ClientCode = req.ClientCode
		if err := tx.FrpClientCfg.Save(cfg); err != nil {
			log.Error("迁移配置隧道失败", zap.Error(err))
			return errors.New("操作失败")
		}
		if cfg.Enable == 1 {
			return engine.ClientCfgConfig(tx, cfg.Code)
		}
		return nil
	})
}
