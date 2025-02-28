package service

import (
	"server/model"
	"server/pkg/bean"
	"server/pkg/jwt"
	"server/repository"
)

type PageReq struct {
	bean.PageParam
	Level      string `json:"level"`
	ClientCode string `binding:"required" json:"clientCode"`
}

type Item struct {
	Id        int    `json:"id"`
	Level     string `json:"level"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"createdAt"`
}

func (service *service) Page(claims jwt.Claims, req PageReq) (list []Item, total int64) {
	db, _, _ := repository.Get("")
	if db.Where("code = ? AND user_code = ?", req.ClientCode, claims.Code).First(&model.GostClient{}).RowsAffected == 0 {
		return nil, 0
	}
	var loggers []model.GostClientLogger
	var where = db.Where("client_code = ?", req.ClientCode)
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
