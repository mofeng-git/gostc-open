package service

import (
	"gorm.io/gen"
	"server/pkg/bean"
	"server/repository"
)

type PageReq struct {
	bean.PageParam
	Level    string `json:"level"`
	NodeCode string `binding:"required" json:"nodeCode"`
}

type Item struct {
	Id        int    `json:"id"`
	Level     string `json:"level"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"createdAt"`
}

func (service *service) Page(req PageReq) (list []Item, total int64) {
	db, _, _ := repository.Get("")
	var where = []gen.Condition{
		db.GostNodeLogger.NodeCode.Eq(req.NodeCode),
	}
	if req.Level != "" {
		where = append(where, db.GostNodeLogger.Level.Eq(req.Level))
	}
	loggers, total, _ := db.GostNodeLogger.Where(where...).Order(db.GostNodeLogger.Id.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	for _, logger := range loggers {
		list = append(list, Item{
			Id:        logger.Id,
			Level:     logger.Level,
			Content:   logger.Content,
			CreatedAt: logger.CreatedAt,
		})
	}
	return list, total
}
