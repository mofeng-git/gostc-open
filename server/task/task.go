package task

import (
	"server/bootstrap"
	"server/global"
)

func init() {
	bootstrap.TaskFunc = func() {
		_, _ = global.Cron.AddFunc("0 0 * * *", gostObs)
	}
}
