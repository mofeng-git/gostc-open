package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/model"
	"server/repository"
	"server/service/common/cache"
)

type BaseReq struct {
	Title   string `binding:"required" json:"title" label:"标题"`
	Favicon string `binding:"required" json:"favicon" label:"图标"`
	BaseUrl string `binding:"required" json:"baseUrl"`
}

func (service *service) Base(req BaseReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *gorm.DB) error {
		tx.Where("kind = ?", model.SYSTEM_CONFIG_KIND_BASE).Delete(&model.SystemConfig{})
		if err := tx.Create(model.GenerateSystemConfigBase(req.Title, req.Favicon, req.BaseUrl)).Error; err != nil {
			log.Error("修改系统基础配置失败", zap.Error(err))
			return errors.New("操作失败")
		}
		cache.SetSystemConfigBase(model.SystemConfigBase{
			Title:   req.Title,
			Favicon: req.Favicon,
			BaseUrl: req.BaseUrl,
		})
		return nil
	})
}
