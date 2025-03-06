package service

import (
	"errors"
	"go.uber.org/zap"
	"server/model"
	"server/repository"
	"server/repository/query"
	"server/service/common/cache"
)

type BaseReq struct {
	Title   string `binding:"required" json:"title" label:"标题"`
	Favicon string `binding:"required" json:"favicon" label:"图标"`
	BaseUrl string `binding:"required" json:"baseUrl"`
	ApiKey  string `json:"apiKey"`
}

func (service *service) Base(req BaseReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		_, _ = tx.SystemConfig.Where(tx.SystemConfig.Kind.Eq(model.SYSTEM_CONFIG_KIND_BASE)).Delete()
		if err := tx.SystemConfig.Create(model.GenerateSystemConfigBase(req.Title, req.Favicon, req.BaseUrl, req.ApiKey)...); err != nil {
			log.Error("修改系统基础配置失败", zap.Error(err))
			return errors.New("操作失败")
		}
		cache.SetSystemConfigBase(model.SystemConfigBase{
			Title:   req.Title,
			Favicon: req.Favicon,
			BaseUrl: req.BaseUrl,
			ApiKey:  req.ApiKey,
		})
		return nil
	})
}
