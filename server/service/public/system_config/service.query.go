package service

import (
	"server/global"
	"server/model"
	"server/repository/cache"
)

type QueryResp struct {
	Title        string `json:"title"`
	Favicon      string `json:"favicon"`
	BaseUrl      string `json:"baseUrl"`
	Version      string `json:"version"`
	Register     string `json:"register"`
	CheckIn      string `json:"checkIn"`
	CheckInStart int    `json:"checkInStart"`
	CheckInEnd   int    `json:"checkInEnd"`

	// 功能菜单
	FuncWeb     string `json:"funcWeb"`
	FuncForward string `json:"funcForward"`
	FuncTunnel  string `json:"funcTunnel"`
	FuncP2P     string `json:"funcP2P"`
	FuncProxy   string `json:"funcProxy"`
	FuncNode    string `json:"funcNode"`

	System QuerySystemInfo `json:"system"`
}

type QuerySystemInfo struct {
	Version string `json:"version"`
}

func (service *service) Query() QueryResp {
	var baseConfig model.SystemConfigBase
	cache.GetSystemConfigBase(&baseConfig)
	var gostConfig model.SystemConfigGost
	cache.GetSystemConfigGost(&gostConfig)
	return QueryResp{
		Title:        baseConfig.Title,
		Favicon:      baseConfig.Favicon,
		BaseUrl:      baseConfig.BaseUrl,
		Version:      gostConfig.Version,
		Register:     baseConfig.Register,
		CheckIn:      baseConfig.CheckIn,
		CheckInStart: baseConfig.CheckInStart,
		CheckInEnd:   baseConfig.CheckInEnd,

		FuncWeb:     gostConfig.FuncWeb,
		FuncForward: gostConfig.FuncForward,
		FuncTunnel:  gostConfig.FuncTunnel,
		FuncP2P:     gostConfig.FuncP2P,
		FuncProxy:   gostConfig.FuncProxy,
		FuncNode:    gostConfig.FuncNode,

		System: QuerySystemInfo{
			Version: global.VERSION,
		},
	}
}
