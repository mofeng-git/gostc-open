package service

import (
	"server/model"
	"server/repository"
	"time"
)

type ListReq struct {
	Used     int    `json:"used"`
	Account  string `json:"account"`
	UserCode string `json:"userCode"`
}

type ListItem struct {
	Code         string   `json:"code"`
	Name         string   `json:"name"`
	UserCode     string   `json:"userCode"`
	ChargingType int      `json:"chargingType"`
	Cycle        int      `json:"cycle"`
	Amount       string   `json:"amount"`
	Limiter      int      `json:"limiter"`
	RLimiter     int      `json:"rLimiter"`
	CLimiter     int      `json:"cLimiter"`
	OnlyChina    int      `json:"onlyChina"`
	Nodes        []string `json:"nodes"`
	TunnelType   int      `json:"tunnelType"`
	TunnelCode   string   `json:"tunnelCode"`
	ExpAt        string   `json:"expAt"`
}

func (service *service) List(req ListReq) (list []ListItem) {
	db, _, _ := repository.Get("")
	var cfgs []model.GostClientConfig
	var where = db
	if req.Account != "" {
		where = where.Where(
			"user_code in (?)",
			db.Model(&model.SystemUser{}).Where("account like ?", "%"+req.Account+"%").Select("code"),
		)
	}
	if req.Used == 1 {
		where = where.Where("tunnel_code != ?", "")
	}
	if req.Used == 2 {
		where = where.Where("tunnel_code == ?", "")
	}
	if req.UserCode != "" {
		where = where.Where("user_code = ?", req.UserCode)
	}
	db.Where(where).Order("id desc").Find(&cfgs)
	for _, cfg := range cfgs {
		list = append(list, ListItem{
			ChargingType: cfg.ChargingType,
			Cycle:        cfg.Cycle,
			Amount:       cfg.Amount.String(),
			Limiter:      cfg.Limiter,
			RLimiter:     cfg.RLimiter,
			CLimiter:     cfg.CLimiter,
			OnlyChina:    cfg.OnlyChina,
			ExpAt:        time.Unix(cfg.ExpAt, 0).Format(time.DateTime),
		})
	}
	return list
}
