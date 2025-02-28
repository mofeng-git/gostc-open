package service

import (
	"gorm.io/gorm"
	"server/model"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/service/common/cache"
	"server/service/common/node_rule"
)

type ListReq struct {
	Bind int `json:"bind"`
}

type ListItem struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Remark string `json:"remark"`

	Web     int `json:"web"`
	Tunnel  int `json:"tunnel"`
	Forward int `json:"forward"`

	Rules     []string         `json:"rules"`
	RuleNames []string         `json:"ruleNames"`
	Tags      []string         `json:"tags"`
	Configs   []ListItemConfig `json:"configs"`
	Online    int              `json:"online"`
}

type ListItemConfig struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	ChargingType int    `json:"chargingType"`
	Cycle        int    `json:"cycle"`
	Amount       string `json:"amount"`
	Limiter      int    `json:"limiter"`
	RLimiter     int    `json:"rLimiter"`
	CLimiter     int    `json:"cLimiter"`
	OnlyChina    int    `json:"onlyChina"`
}

func (service *service) List(claims jwt.Claims, req ListReq) (list []ListItem) {
	db, _, _ := repository.Get("")
	var nodes []model.GostNode
	var excludeCodes []string
	db.Model(&model.GostNodeBind{}).Where("user_code != ?", claims.Code).Pluck("node_code", &excludeCodes)

	var where = db
	if len(excludeCodes) > 0 {
		where = where.Where("code NOT IN ?", excludeCodes)
	}
	switch req.Bind {
	case 1:
		where = where.Where("code IN (?)",
			db.Model(&model.GostNodeBind{}).Where("user_code = ?", claims.Code).Select("node_code"),
		)
	case 2:
		where = where.Where("code NOT IN (?)",
			db.Model(&model.GostNodeBind{}).Where("user_code = ?", claims.Code).Select("node_code"),
		)
	}
	db.Preload("Configs", func(db *gorm.DB) *gorm.DB {
		return db.Order("gost_node_config.index_value asc")
	}).Where(where).Order("index_value asc").Order("id desc").Find(&nodes)
	for _, node := range nodes {
		if len(node.Configs) == 0 {
			continue
		}
		var configs []ListItemConfig
		for _, cfg := range node.Configs {
			configs = append(configs, ListItemConfig{
				Code:         cfg.Code,
				Name:         cfg.Name,
				ChargingType: cfg.ChargingType,
				Cycle:        cfg.Cycle,
				Amount:       cfg.Amount.String(),
				Limiter:      cfg.Limiter,
				RLimiter:     cfg.RLimiter,
				CLimiter:     cfg.CLimiter,
				OnlyChina:    cfg.OnlyChina,
			})
		}
		var ruleNames []string
		for _, rule := range node.GetRules() {
			ruleNames = append(ruleNames, node_rule.RuleMap[rule].Name())
		}
		list = append(list, ListItem{
			Code:      node.Code,
			Name:      node.Name,
			Remark:    node.Remark,
			Web:       node.Web,
			Tunnel:    node.Tunnel,
			Forward:   node.Forward,
			Rules:     node.GetRules(),
			RuleNames: ruleNames,
			Tags:      node.GetTags(),
			Configs:   configs,
			Online:    utils.TrinaryOperation(cache.GetNodeOnline(node.Code), 1, 2),
		})
	}
	return list
}
