package service

import (
	"server/model"
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
	OnlyChina    int    `json:"onlyChina"`
	IndexValue   int    `json:"indexValue"`
}

func (service *service) List(req ListReq) (list []ListItem) {
	db, _, _ := repository.Get("")
	var cfgs []model.GostNodeConfig
	var where = db.Where("node_code = ?", req.NodeCode)
	db.Where(where).Order("index_value asc").Order("id desc").Find(&cfgs)
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
			OnlyChina:    cfg.OnlyChina,
			IndexValue:   cfg.IndexValue,
		})
	}
	return list
}
