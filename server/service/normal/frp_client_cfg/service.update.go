package service

import (
	"errors"
	"go.uber.org/zap"
	"server/pkg/jwt"
	"server/repository"
	"server/repository/query"
	"server/service/engine"
)

type UpdateReq struct {
	Code     string `binding:"required" json:"code"`
	Name     string `binding:"required" json:"name"`
	Type     string `binding:"required" json:"type"`
	Content  string `binding:"required" json:"content"`
	Address  string `json:"address"`
	Platform string `json:"platform"`
}

func (service *service) Update(claims jwt.Claims, req UpdateReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(claims.Code)).First()
		if user == nil {
			return errors.New("用户错误")
		}

		cfg, _ := tx.FrpClientCfg.Where(
			tx.FrpClientCfg.UserCode.Eq(user.Code),
			tx.FrpClientCfg.Code.Eq(req.Code),
		).First()
		if cfg == nil {
			return errors.New("操作失败")
		}

		cfg.Name = req.Name
		cfg.ContentType = req.Type
		cfg.Content = req.Content
		cfg.Platform = req.Platform
		cfg.Address = req.Address

		if err := tx.FrpClientCfg.Save(cfg); err != nil {
			log.Error("修改配置隧道失败", zap.Error(err))
			return errors.New("操作失败")
		}

		if cfg.Enable == 1 {
			return engine.ClientCfgConfig(tx, cfg.Code)
		}
		return nil
	})
}
