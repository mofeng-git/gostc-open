package service

import (
	"gorm.io/gen"
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
	client, _ := db.GostClient.Where(db.GostClient.UserCode.Eq(claims.Code), db.GostClient.Code.Eq(req.ClientCode)).First()
	if client == nil {
		return nil, 0
	}
	var where = []gen.Condition{
		db.GostClientLogger.ClientCode.Eq(req.ClientCode),
	}
	if req.Level != "" {
		where = append(where, db.GostClientLogger.Level.Eq(req.Level))
	}

	logs, total, _ := db.GostClientLogger.Where(where...).Order(db.GostClientLogger.Id.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	for _, logger := range logs {
		list = append(list, Item{
			Id:        logger.Id,
			Level:     logger.Level,
			Content:   logger.Content,
			CreatedAt: logger.CreatedAt,
		})
	}
	return list, total
}
