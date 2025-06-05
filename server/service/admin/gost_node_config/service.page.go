package service

import (
	"gorm.io/gen"
	"server/pkg/bean"
	"server/repository"
)

type PageReq struct {
	bean.PageParam
	Name string `json:"name"`
}

type Item struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	ChargingType int    `json:"chargingType"`
	Cycle        int    `json:"cycle"`
	Amount       string `json:"amount"`
	Limiter      int    `json:"limiter"`
	RLimiter     int    `json:"rLimiter"`
	CLimiter     int    `json:"cLimiter"`
	IndexValue   int    `json:"indexValue"`
}

func (service *service) Page(req PageReq) (list []Item, total int64) {
	db, _, _ := repository.Get("")
	var where []gen.Condition
	if req.Name != "" {
		where = append(where, db.GostNodeConfig.Name.Like("%"+req.Name+"%"))
	}

	cfgs, total, _ := db.GostNodeConfig.Where(where...).Order(
		db.GostNodeConfig.IndexValue.Asc(),
		db.GostNodeConfig.Id.Desc(),
	).FindByPage(req.GetOffset(), req.GetLimit())
	for _, cfg := range cfgs {
		list = append(list, Item{
			Code:         cfg.Code,
			Name:         cfg.Name,
			ChargingType: cfg.ChargingType,
			Cycle:        cfg.Cycle,
			Amount:       cfg.Amount.String(),
			Limiter:      cfg.Limiter,
			RLimiter:     cfg.RLimiter,
			CLimiter:     cfg.CLimiter,
			IndexValue:   cfg.IndexValue,
		})
	}
	return list, total
}
