package task

import (
	"server/bootstrap"
	"server/model"
	"server/repository"
	"time"
)

func init() {
	bootstrap.AddCron("0 0 * * *", func() {
		db, _, _ := repository.Get("")
		db.Where("created_at <= ?", time.Now().AddDate(0, 0, -3)).Delete(&model.GostNodeLogger{})
	})
}
