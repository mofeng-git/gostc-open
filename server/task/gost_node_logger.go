package task

import (
	"server/repository"
	"time"
)

func gostNodeLogger() {
	db, _, _ := repository.Get("")
	_, _ = db.GostNodeLogger.Where(db.GostNodeLogger.CreatedAt.Lte(time.Now().AddDate(0, 0, -3).Unix())).Delete()
}
