package todo

import (
	"server/model"
	"server/repository"
	"server/service/common/cache"
)

func systemConfig() {
	db, _, _ := repository.Get("")
	for _, item := range model.GenerateSystemConfigBase(
		"GOSTC",
		"https://oss.sian.one/static/master/325ebdbe993146e08f9a7aa5b8e59d02.png",
		"https://gost.sian.one",
		"",
		"2",
		"2",
		1,
		5,
	) {
		_ = db.SystemConfig.Create(item)
	}
	for _, item := range model.GenerateSystemConfigGost(
		"v1.0.0",
		"2",
		"1",
		"1",
		"1",
		"1",
		"1",
		"1",
		"1",
	) {
		_ = db.SystemConfig.Create(item)
	}

	configs, _ := db.SystemConfig.Find()
	baseConfig := model.GetSystemConfigBase(configs)
	cache.SetSystemConfigBase(baseConfig)
	gostConfig := model.GetSystemConfigGost(configs)
	cache.SetSystemConfigGost(gostConfig)
}
