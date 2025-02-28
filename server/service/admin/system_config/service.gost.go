package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/model"
	"server/repository"
	"server/service/common/cache"
)

type GostReq struct {
	Version string `binding:"required" json:"version"`
	Logger  string `binding:"required" json:"logger"`
}

func (service *service) Gost(req GostReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *gorm.DB) error {
		tx.Where("kind = ?", model.SYSTEM_CONFIG_KIND_GOST).Delete(&model.SystemConfig{})
		if err := tx.Create(model.GenerateSystemConfigGost(req.Version, req.Logger)).Error; err != nil {
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
