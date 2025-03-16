package service

import (
	"server/model"
	"server/service/common/cache"
)

type QueryResp struct {
	Title    string `json:"title"`
	Favicon  string `json:"favicon"`
	BaseUrl  string `json:"baseUrl"`
	Version  string `json:"version"`
	Register string `json:"register"`
	CheckIn  string `json:"checkIn"`
}

func (service *service) Query() QueryResp {
	var baseConfig model.SystemConfigBase
	cache.GetSystemConfigBase(&baseConfig)
	var gostConfig model.SystemConfigGost
	cache.GetSystemConfigGost(&gostConfig)
	return QueryResp{
		Title:    baseConfig.Title,
		Favicon:  baseConfig.Favicon,
		BaseUrl:  baseConfig.BaseUrl,
		Version:  gostConfig.Version,
		Register: baseConfig.Register,
		CheckIn:  baseConfig.CheckIn,
	}
}
