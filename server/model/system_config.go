package model

const (
	SYSTEM_CONFIG_KIND_BASE = "SystemConfigBase" // 基础配置
	SYSTEM_CONFIG_KIND_GOST = "SystemConfigGost" // Gost配置
)

type SystemConfig struct {
	Kind  string `gorm:"column:kind;uniqueIndex:system_config_uidx"`
	Key   string `gorm:"column:key;uniqueIndex:system_config_uidx"`
	Value string `gorm:"column:value"`
}

type SystemConfigBase struct {
	Title   string `json:"title"`
	Favicon string `json:"favicon"`
	BaseUrl string `json:"baseUrl"`
}

func GetSystemConfigBase(list []SystemConfig) (result SystemConfigBase) {
	for _, item := range list {
		switch item.Key {
		case "title":
			result.Title = item.Value
		case "favicon":
			result.Favicon = item.Value
		case "baseUrl":
			result.BaseUrl = item.Value
		}
	}
	return result
}

func GenerateSystemConfigBase(title, favicon, baseUrl string) []SystemConfig {
	return []SystemConfig{
		{Kind: SYSTEM_CONFIG_KIND_BASE, Key: "title", Value: title},
		{Kind: SYSTEM_CONFIG_KIND_BASE, Key: "favicon", Value: favicon},
		{Kind: SYSTEM_CONFIG_KIND_BASE, Key: "baseUrl", Value: baseUrl},
	}
}

type SystemConfigGost struct {
	Version string `json:"version"`
	Logger  string `json:"logger"`
}

func GenerateSystemConfigGost(version string, logger string) []SystemConfig {
	return []SystemConfig{
		{Kind: SYSTEM_CONFIG_KIND_GOST, Key: "version", Value: version},
		{Kind: SYSTEM_CONFIG_KIND_GOST, Key: "logger", Value: logger},
	}
}

func GetSystemConfigGost(list []SystemConfig) (result SystemConfigGost) {
	for _, item := range list {
		switch item.Key {
		case "version":
			result.Version = item.Value
		case "logger":
			result.Logger = item.Value
		}
	}
	return result
}
