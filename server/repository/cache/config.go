package cache

import (
	"github.com/patrickmn/go-cache"
	"server/global"
	"server/model"
)

const (
    system_config_base_key  = "system:config:base"
    system_config_gost_key  = "system:config:gost"
    system_config_email_key = "system:config:email"
    system_config_home_key  = "system:config:home"
)

func SetSystemConfigBase(cfg model.SystemConfigBase) {
	global.Cache.SetStruct(system_config_base_key, cfg, cache.NoExpiration)
}

func GetSystemConfigBase(cfg *model.SystemConfigBase) {
	_ = global.Cache.GetStruct(system_config_base_key, cfg)
}

func SetSystemConfigGost(cfg model.SystemConfigGost) {
	global.Cache.SetStruct(system_config_gost_key, cfg, cache.NoExpiration)
}

func GetSystemConfigGost(cfg *model.SystemConfigGost) {
	_ = global.Cache.GetStruct(system_config_gost_key, cfg)
}

func SetSystemConfigEmail(cfg model.SystemConfigEmail) {
	global.Cache.SetStruct(system_config_email_key, cfg, cache.NoExpiration)
}

func GetSystemConfigEmail(cfg *model.SystemConfigEmail) {
    _ = global.Cache.GetStruct(system_config_email_key, cfg)
}

func SetSystemConfigHome(cfg model.SystemConfigHome) {
    global.Cache.SetStruct(system_config_home_key, cfg, cache.NoExpiration)
}

func GetSystemConfigHome(cfg *model.SystemConfigHome) {
    _ = global.Cache.GetStruct(system_config_home_key, cfg)
}
