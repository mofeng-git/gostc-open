package bootstrap

import (
	"server/global"
	"server/pkg/memory/local"
	"time"
)

func InitMemory() {
	global.Cache = local.NewMemory(global.BASE_PATH+"/data/cache.db", time.Minute*10)
	releaseFunc = append(releaseFunc, func() {
		global.Cache.Sync()
	})
}
