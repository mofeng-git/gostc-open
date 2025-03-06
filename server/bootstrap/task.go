package bootstrap

import (
	"github.com/robfig/cron/v3"
	"server/global"
)

var TaskFunc func()

func InitTask() {
	global.Cron = cron.New()
	if TaskFunc != nil {
		TaskFunc()
	}
	global.Cron.Start()
}
