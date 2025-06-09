package service

import (
	"errors"
	"go.uber.org/zap"
	"server/pkg/jwt"
	"server/repository"
	"server/repository/query"
	"server/service/engine"
)

type DeleteReq struct {
	Code string `binding:"required" json:"code" label:"编号"`
}

func (service *service) Delete(claims jwt.Claims, req DeleteReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		cfg, _ := tx.FrpClientCfg.Where(
			tx.FrpClientCfg.UserCode.Eq(claims.Code),
			tx.FrpClientCfg.Code.Eq(req.Code),
		).First()
		if cfg == nil {
			return errors.New("操作失败")
		}
		if _, err := tx.FrpClientCfg.Where(tx.FrpClientCfg.Code.Eq(cfg.Code)).Delete(); err != nil {
			log.Error("删除用户配置隧道失败", zap.Error(err))
			return errors.New("操作失败")
		}
		engine.ClientCfgRemove(cfg.ClientCode, cfg.Code)
		return nil
	})
}
