package service

import (
	"server/model"
	"server/repository"
	"server/service/common/cache"
)

type QueryReq struct {
	Kind string `binding:"required" json:"kind"`
}

func (service *service) Query(req QueryReq) any {
	db, _, _ := repository.Get("")
	var configs []model.SystemConfig
	db.Where("kind = ?", req.Kind).Find(&configs)
	switch req.Kind {
	case model.SYSTEM_CONFIG_KIND_BASE:
		var cfg model.SystemConfigBase
		cache.GetSystemConfigBase(&cfg)
		return cfg
	case model.SYSTEM_CONFIG_KIND_GOST:
		var cfg model.SystemConfigGost
		cache.GetSystemConfigGost(&cfg)
		return cfg
	}
	return nil
}
