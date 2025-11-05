package todo

import (
	"server/model"
	"server/repository"
	"server/repository/cache"
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
        "1",
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

    for _, item := range model.GenerateSystemConfigEmail(
        "2",
        "管理员",
        "",
        "",
        "",
        "",
        "您正在进行重置密码操作，请勿将验证码告诉他人，验证码：{{CODE}}(5分钟有效)，发件时间：{{DATETIME}}",
    ) {
        _ = db.SystemConfig.Create(item)
    }

    // 默认首页配置：开启+使用内置模板
    for _, item := range model.GenerateSystemConfigHome(
        "1", // 开启
        "",  // 留空则使用内置模板
    ) {
        _ = db.SystemConfig.Create(item)
    }

	configs, _ := db.SystemConfig.Find()
    baseConfig := model.GetSystemConfigBase(configs)
    cache.SetSystemConfigBase(baseConfig)
    gostConfig := model.GetSystemConfigGost(configs)
    cache.SetSystemConfigGost(gostConfig)
    emailConfig := model.GetSystemConfigEmail(configs)
    cache.SetSystemConfigEmail(emailConfig)
    homeConfig := model.GetSystemConfigHome(configs)
    cache.SetSystemConfigHome(homeConfig)
}
