package service

import (
	"errors"
	"go.uber.org/zap"
	"server/model"
	"server/repository"
	"server/repository/query"
	"server/service/common/cache"
)

type GostReq struct {
	Version string `binding:"required" json:"version"`
	Logger  string `binding:"required" json:"logger"`
}

func (service *service) Gost(req GostReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		_, _ = tx.SystemConfig.Where(tx.SystemConfig.Kind.Eq(model.SYSTEM_CONFIG_KIND_GOST)).Delete()
		if err := tx.SystemConfig.Create(model.GenerateSystemConfigGost(req.Version, req.Logger)...); err != nil {
			log.Error("修改系统GOST配置失败", zap.Error(err))
			return errors.New("操作失败")
		}
		cache.SetSystemConfigGost(model.SystemConfigGost{
			Version: req.Version,
			Logger:  req.Logger,
		})
		return nil
	})
}
