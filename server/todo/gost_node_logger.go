package todo

import (
	"server/bootstrap"
	"server/model"
	"server/repository"
	"time"
)

func init() {
	bootstrap.AddTodo(func() {
		db, _, _ := repository.Get("")
		db.Where("created_at <= ?", time.Now().AddDate(0, 0, -3).Unix()).Delete(&model.GostNodeLogger{})
	})
}
