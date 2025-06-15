package service

import (
	"errors"
	"go.uber.org/zap"
	"server/model"
	"server/repository"
	"server/repository/cache"
	"server/repository/query"
	"strings"
)

type BaseReq struct {
	Title        string `binding:"required" json:"title" label:"标题"`
	Favicon      string `binding:"required" json:"favicon" label:"图标"`
	BaseUrl      string `binding:"required" json:"baseUrl"`
	ApiKey       string `json:"apiKey"`
	Register     string `binding:"required" json:"register"`
	CheckIn      string `binding:"required" json:"checkIn"`
	CheckInStart int    `json:"checkInStart"`
	CheckInEnd   int    `json:"checkInEnd"`
}

func (service *service) Base(req BaseReq) error {
	if !strings.HasPrefix(req.BaseUrl, "http://") && !strings.HasPrefix(req.BaseUrl, "https://") {
		return errors.New("基础URL必须http://或https://开头")
	}
	req.BaseUrl = strings.TrimRight(req.BaseUrl, "/")
	db, _, log := repository.Get("")
	if req.CheckInStart > req.CheckInEnd {
		req.CheckInStart, req.CheckInEnd = req.CheckInEnd, req.CheckInStart
	}
	return db.Transaction(func(tx *query.Query) error {
		_, _ = tx.SystemConfig.Where(tx.SystemConfig.Kind.Eq(model.SYSTEM_CONFIG_KIND_BASE)).Delete()
		if err := tx.SystemConfig.Create(model.GenerateSystemConfigBase(
			req.Title,
			req.Favicon,
			req.BaseUrl,
			req.ApiKey,
			req.Register,
			req.CheckIn,
			req.CheckInStart,
			req.CheckInEnd,
		)...); err != nil {
			log.Error("修改系统基础配置失败", zap.Error(err))
			return errors.New("操作失败")
		}
		cache.SetSystemConfigBase(model.SystemConfigBase{
			Title:        req.Title,
			Favicon:      req.Favicon,
			BaseUrl:      req.BaseUrl,
			ApiKey:       req.ApiKey,
			Register:     req.Register,
			CheckIn:      req.CheckIn,
			CheckInStart: req.CheckInStart,
			CheckInEnd:   req.CheckInEnd,
		})
		return nil
	})
}
