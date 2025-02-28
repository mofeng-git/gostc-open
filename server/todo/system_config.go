package todo

import (
	"server/bootstrap"
	"server/model"
	"server/repository"
	"server/service/common/cache"
)

func init() {
	bootstrap.AddTodo(func() {
		db, _, _ := repository.Get("")
		db.Create(model.GenerateSystemConfigBase("GOSTC", "https://oss.sian.one/static/master/325ebdbe993146e08f9a7aa5b8e59d02.png", "https://gost.sian.one"))
		db.Create(model.GenerateSystemConfigGost("v1.0.0", "2"))

		var configs []model.SystemConfig
		db.Find(&configs)
		baseConfig := model.GetSystemConfigBase(configs)
		cache.SetSystemConfigBase(baseConfig)
		gostConfig := model.GetSystemConfigGost(configs)
		cache.SetSystemConfigGost(gostConfig)
	})
}
