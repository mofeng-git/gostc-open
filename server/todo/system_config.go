package todo

import (
	"server/model"
	"server/repository"
	"server/service/common/cache"
)

func systemConfig() {
	db, _, _ := repository.Get("")
	_ = db.SystemConfig.Create(model.GenerateSystemConfigBase("GOSTC", "https://oss.sian.one/static/master/325ebdbe993146e08f9a7aa5b8e59d02.png", "https://gost.sian.one", "")...)
	_ = db.SystemConfig.Create(model.GenerateSystemConfigGost("v1.0.0", "2")...)

	configs, _ := db.SystemConfig.Find()
	baseConfig := model.GetSystemConfigBase(configs)
	cache.SetSystemConfigBase(baseConfig)
	gostConfig := model.GetSystemConfigGost(configs)
	cache.SetSystemConfigGost(gostConfig)
}
