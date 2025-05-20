package model

import "strconv"

const (
	SYSTEM_CONFIG_KIND_BASE = "SystemConfigBase" // 基础配置
	SYSTEM_CONFIG_KIND_GOST = "SystemConfigGost" // Gost配置
)

type SystemConfig struct {
	Kind  string `gorm:"column:kind;size:100;uniqueIndex:system_config_uidx"`
	Key   string `gorm:"column:key;size:100;uniqueIndex:system_config_uidx"`
	Value string `gorm:"column:value"`
}

type SystemConfigBase struct {
	Title        string `json:"title"`
	Favicon      string `json:"favicon"`
	BaseUrl      string `json:"baseUrl"`
	ApiKey       string `json:"apiKey"`
	Register     string `json:"register"`
	CheckIn      string `json:"checkIn"`
	CheckInStart int    `json:"checkInStart"`
	CheckInEnd   int    `json:"checkInEnd"`
}

func GetSystemConfigBase(list []*SystemConfig) (result SystemConfigBase) {
	for _, item := range list {
		switch item.Key {
		case "title":
			result.Title = item.Value
		case "favicon":
			result.Favicon = item.Value
		case "baseUrl":
			result.BaseUrl = item.Value
		case "apiKey":
			result.ApiKey = item.Value
		case "register":
			result.Register = item.Value
		case "checkIn":
			result.CheckIn = item.Value
		case "checkInStart":
			result.CheckInStart, _ = strconv.Atoi(item.Value)
		case "checkInEnd":
			result.CheckInEnd, _ = strconv.Atoi(item.Value)
		}
	}
	return result
}

func GenerateSystemConfigBase(title, favicon, baseUrl, apiKey, register, checkIn string, checkInStart, checkInEnd int) []*SystemConfig {
	return []*SystemConfig{
		{Kind: SYSTEM_CONFIG_KIND_BASE, Key: "title", Value: title},
		{Kind: SYSTEM_CONFIG_KIND_BASE, Key: "favicon", Value: favicon},
		{Kind: SYSTEM_CONFIG_KIND_BASE, Key: "baseUrl", Value: baseUrl},
		{Kind: SYSTEM_CONFIG_KIND_BASE, Key: "apiKey", Value: apiKey},
		{Kind: SYSTEM_CONFIG_KIND_BASE, Key: "register", Value: register},
		{Kind: SYSTEM_CONFIG_KIND_BASE, Key: "checkIn", Value: checkIn},
		{Kind: SYSTEM_CONFIG_KIND_BASE, Key: "checkInStart", Value: strconv.Itoa(checkInStart)},
		{Kind: SYSTEM_CONFIG_KIND_BASE, Key: "checkInEnd", Value: strconv.Itoa(checkInEnd)},
	}
}

type SystemConfigGost struct {
	Version string `json:"version"`
	Logger  string `json:"logger"`

	// 功能菜单
	FuncWeb     string `json:"funcWeb"`
	FuncForward string `json:"funcForward"`
	FuncTunnel  string `json:"funcTunnel"`
	FuncP2P     string `json:"funcP2P"`
	FuncProxy   string `json:"funcProxy"`
	FuncTun     string `json:"funcTun"`
	FuncNode    string `json:"funcNode"`
}

func GenerateSystemConfigGost(version string, logger string, funcWeb, funcForward, funcTunnel, funcP2P, funcProxy, funcTun, funcNode string) []*SystemConfig {
	return []*SystemConfig{
		{Kind: SYSTEM_CONFIG_KIND_GOST, Key: "version", Value: version},
		{Kind: SYSTEM_CONFIG_KIND_GOST, Key: "logger", Value: logger},
		{Kind: SYSTEM_CONFIG_KIND_GOST, Key: "funcWeb", Value: funcWeb},
		{Kind: SYSTEM_CONFIG_KIND_GOST, Key: "funcForward", Value: funcForward},
		{Kind: SYSTEM_CONFIG_KIND_GOST, Key: "funcTunnel", Value: funcTunnel},
		{Kind: SYSTEM_CONFIG_KIND_GOST, Key: "funcP2P", Value: funcP2P},
		{Kind: SYSTEM_CONFIG_KIND_GOST, Key: "funcProxy", Value: funcProxy},
		{Kind: SYSTEM_CONFIG_KIND_GOST, Key: "funcTun", Value: funcTun},
		{Kind: SYSTEM_CONFIG_KIND_GOST, Key: "funcNode", Value: funcNode},
	}
}

func GetSystemConfigGost(list []*SystemConfig) (result SystemConfigGost) {
	for _, item := range list {
		switch item.Key {
		case "version":
			result.Version = item.Value
		case "logger":
			result.Logger = item.Value
		case "funcWeb":
			result.FuncWeb = item.Value
		case "funcForward":
			result.FuncForward = item.Value
		case "funcTunnel":
			result.FuncTunnel = item.Value
		case "funcP2P":
			result.FuncP2P = item.Value
		case "funcProxy":
			result.FuncProxy = item.Value
		case "funcTun":
			result.FuncTun = item.Value
		case "funcNode":
			result.FuncNode = item.Value
		}
	}
	return result
}
