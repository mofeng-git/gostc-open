package service

import (
	"gorm.io/gen"
	"server/pkg/bean"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/repository/cache"
)

type PageReq struct {
	bean.PageParam
	Name       string `json:"name"`
	Enable     int    `json:"enable"`
	ClientCode string `json:"clientCode"`
}

type Item struct {
	Code     string     `json:"code"`
	Name     string     `json:"name"`
	Type     string     `json:"type"`
	Content  string     `json:"content"`
	Address  string     `json:"address"`
	Platform string     `json:"platform"`
	Enable   int        `json:"enable"`
	Client   ItemClient `json:"client"`
}

type ItemClient struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Online int    `json:"online"`
}

func (service *service) Page(claims jwt.Claims, req PageReq) (list []Item, total int64) {
	db, _, _ := repository.Get("")
	var where = []gen.Condition{
		db.FrpClientCfg.UserCode.Eq(claims.Code),
	}
	if req.Name != "" {
		where = append(where, db.FrpClientCfg.Name.Like("%"+req.Name+"%"))
	}
	if req.Enable > 0 {
		where = append(where, db.FrpClientCfg.Enable.Eq(req.Enable))
	}
	if req.ClientCode != "" {
		where = append(where, db.FrpClientCfg.ClientCode.Eq(req.ClientCode))
	}
	cfgs, total, _ := db.FrpClientCfg.Preload(
		db.FrpClientCfg.User,
		db.FrpClientCfg.Client,
	).Where(where...).Order(db.FrpClientCfg.Id.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	for _, cfg := range cfgs {
		list = append(list, Item{
			Code:     cfg.Code,
			Name:     cfg.Name,
			Type:     cfg.ContentType,
			Content:  cfg.Content,
			Address:  cfg.Address,
			Platform: cfg.Platform,
			Enable:   cfg.Enable,
			Client: ItemClient{
				Code:   cfg.ClientCode,
				Name:   cfg.Client.Name,
				Online: utils.TrinaryOperation(cache.GetClientOnline(cfg.ClientCode), 1, 2),
			},
		})
	}
	return list, total
}
