package service

import (
	"gorm.io/gen"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/repository/cache"
	"server/service/common/node_rule"
)

type ListReq struct {
	Bind int `json:"bind"`
}

type ListItem struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Remark string `json:"remark"`

	Web          int    `json:"web"`
	Tunnel       int    `json:"tunnel"`
	Forward      int    `json:"forward"`
	Proxy        int    `json:"proxy"`
	P2P          int    `json:"p2p"`
	ForwardPorts string `json:"forwardPorts"`

	Rules        []string         `json:"rules"`
	RuleNames    []string         `json:"ruleNames"`
	Tags         []string         `json:"tags"`
	Configs      []ListItemConfig `json:"configs"`
	Online       int              `json:"online"`
	CustomDomain int              `json:"customDomain"`
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
}

func (service *service) List(claims jwt.Claims, req ListReq) (list []ListItem) {
	db, _, _ := repository.Get("")

	var excludeCodes []string
	_ = db.GostNodeBind.Where(db.GostNodeBind.UserCode.Neq(claims.Code)).Pluck(db.GostNodeBind.NodeCode, &excludeCodes)

	var where = []gen.Condition{
		db.GostNode.Code.NotIn(excludeCodes...),
	}

	var myNodeCodes []string
	_ = db.GostNodeBind.Where(db.GostNodeBind.UserCode.Eq(claims.Code)).Pluck(db.GostNodeBind.NodeCode, &myNodeCodes)
	switch req.Bind {
	case 1:
		where = append(where, db.GostNode.Code.In(myNodeCodes...))
	case 2:
		where = append(where, db.GostNode.Code.NotIn(myNodeCodes...))
	}
	nodes, _ := db.GostNode.Preload(db.GostNode.Configs.Order(db.GostNodeConfig.IndexValue.Asc())).Where(where...).Order(db.GostNode.IndexValue.Asc(), db.GostNode.Id.Desc()).Find()
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
			})
		}
		var ruleNames []string
		for _, rule := range node.GetRules() {
			getRule := node_rule.Registry.GetRule(rule)
			ruleNames = append(ruleNames, getRule.Name())
		}
		list = append(list, ListItem{
			Code:         node.Code,
			Name:         node.Name,
			Remark:       node.Remark,
			Web:          node.Web,
			Tunnel:       node.Tunnel,
			Forward:      node.Forward,
			Proxy:        node.Proxy,
			P2P:          node.P2P,
			ForwardPorts: node.ForwardPorts,
			Rules:        node.GetRules(),
			RuleNames:    ruleNames,
			Tags:         node.GetTags(),
			Configs:      configs,
			Online:       utils.TrinaryOperation(cache.GetNodeOnline(node.Code), 1, 2),
			CustomDomain: utils.TrinaryOperation(cache.GetNodeCustomDomain(node.Code), 1, 2),
		})
	}
	return list
}
