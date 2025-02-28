package service

import (
	"server/model"
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
	var loggers []model.GostNodeLogger
	var where = db.Where("node_code = ?", req.NodeCode)
	if req.Level != "" {
		where = where.Where("level = ?", req.Level)
	}
	db.Where(where).Model(&loggers).Count(&total)
	db.Where(where).Order("id desc").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&loggers)
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
