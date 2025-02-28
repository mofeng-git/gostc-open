package bootstrap

import (
	"github.com/robfig/cron/v3"
	"server/global"
)

func InitCron() {
	global.Cron = cron.New()
	for _, f := range tasks {
		_, _ = global.Cron.AddFunc(f())
	}
	global.Cron.Start()
	releaseFunc = append(releaseFunc, func() {
		global.Cron.Stop()
	})
}

var tasks []func() (string, func())

func AddCron(spec string, f func()) {
	tasks = append(tasks, func() (string, func()) {
		return spec, f
	})
}
