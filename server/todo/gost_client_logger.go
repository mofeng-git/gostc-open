package todo

import (
	"server/repository"
	"time"
)

func gostClientLogger() {
	db, _, _ := repository.Get("")
	_, _ = db.GostClientLogger.Where(db.GostClientLogger.CreatedAt.Lte(time.Now().AddDate(0, 0, -3).Unix())).Delete()
}
