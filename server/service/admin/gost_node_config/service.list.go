package service

import (
	"server/repository"
)

type ListReq struct {
	NodeCode string `binding:"required" json:"nodeCode"`
}

type ListItem struct {
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

func (service *service) List(req ListReq) (list []ListItem) {
	db, _, _ := repository.Get("")
	cfgs, _ := db.GostNodeConfig.Where(db.GostNodeConfig.NodeCode.Eq(req.NodeCode)).Order(
		db.GostNodeConfig.IndexValue.Asc(),
		db.GostNodeConfig.Id.Desc(),
	).Find()
	for _, cfg := range cfgs {
		list = append(list, ListItem{
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
	return list
}
