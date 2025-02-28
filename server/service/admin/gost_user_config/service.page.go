package service

import (
	"server/pkg/bean"
)

type PageReq struct {
	bean.PageParam
	Account string `json:"account"`
}

type Item struct {
	Code         string   `json:"code"`
	Name         string   `json:"name"`
	UserCode     string   `json:"userCode"`
	UserAccount  string   `json:"userAccount"`
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
	ExpAt        int64    `json:"expAt"`
}

func (service *service) Page(req PageReq) (list []Item, total int64) {
	//var cfgs []model.GostClientConfig
	//var where = db
	//if req.Account != "" {
	//	where = where.Where(
	//		"user_code in (?)",
	//		db.Model(&model.SystemUser{}).Where("account like ?", "%"+req.Account+"%").Select("code"),
	//	)
	//}
	//db.Where(where).Model(&cfgs).Count(&total)
	//db.Preload("User").Where(where).Order("id desc").
	//	Offset(req.GetOffset()).
	//	Limit(req.GetLimit()).
	//	Find(&cfgs)
	//for _, cfg := range cfgs {
	//	list = append(list, Item{
	//		Code:         cfg.Code,
	//		Name:         cfg.Name,
	//		UserCode:     cfg.UserCode,
	//		UserAccount:  cfg.User.Account,
	//		ChargingType: cfg.ChargingType,
	//		Cycle:        cfg.Cycle,
	//		Amount:       cfg.Amount.String(),
	//		Limiter:      cfg.Limiter,
	//		RLimiter:     cfg.RLimiter,
	//		CLimiter:     cfg.CLimiter,
	//		OnlyChina:    cfg.OnlyChina,
	//		Nodes:        cfg.GetNodes(),
	//		TunnelType:   cfg.TunnelType,
	//		TunnelCode:   cfg.TunnelCode,
	//		ExpAt:        cfg.ExpAt,
	//	})
	//}
	return list, total
}
