package service

import (
	"server/model"
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
	OnlyChina    int    `json:"onlyChina"`
	IndexValue   int    `json:"indexValue"`
}

func (service *service) Page(req PageReq) (list []Item, total int64) {
	db, _, _ := repository.Get("")
	var cfgs []model.GostNodeConfig
	var where = db
	if req.Name != "" {
		where = where.Where("name like ?", "%"+req.Name+"%")
	}
	db.Where(where).Model(&cfgs).Count(&total)
	db.Where(where).Order("index_value asc").Order("id desc").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&cfgs)
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
			OnlyChina:    cfg.OnlyChina,
			IndexValue:   cfg.IndexValue,
		})
	}
	return list, total
}
