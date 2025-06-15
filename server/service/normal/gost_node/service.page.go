package service

import (
	"gorm.io/gen"
	"server/model"
	"server/pkg/bean"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	cache2 "server/repository/cache"
	"server/service/common/node_rule"
	"time"
)

type PageReq struct {
	bean.PageParam
}

type Item struct {
	Code   string `json:"code"`
	Key    string `json:"key"`
	Name   string `json:"name"`
	Remark string `json:"remark"`

	Web     int `json:"web"`
	Tunnel  int `json:"tunnel"`
	Forward int `json:"forward"`
	Proxy   int `json:"proxy"`
	P2P     int `json:"p2p"`

	Domain           string `json:"domain"`
	CustomDomain     int    `json:"customDomain"`
	DenyDomainPrefix string `json:"denyDomainPrefix"`
	UrlTpl           string `json:"urlTpl"`
	Address          string `json:"address"`
	ReplaceAddress   string `json:"replaceAddress"`
	Protocol         string `json:"protocol"`

	HttpPort          string   `json:"httpPort"`
	ForwardPorts      string   `json:"forwardPorts"`
	Rules             []string `json:"rules"`
	RuleNames         []string `json:"ruleNames"`
	Tags              []string `json:"tags"`
	IndexValue        int      `json:"indexValue"`
	P2PDisableForward int      `json:"p2pDisableForward"`

	Online      int    `json:"online"`
	Version     string `json:"version"`
	InputBytes  int64  `json:"inputBytes"`
	OutputBytes int64  `json:"outputBytes"`

	LimitResetIndex int   `json:"limitResetIndex"`
	LimitUseTotal   int64 `json:"limitUseTotal"`
	LimitTotal      int   `json:"limitTotal"`
	LimitKind       int   `json:"limitKind"`
}

func (service *service) Page(claims jwt.Claims, req PageReq) (list []Item, total int64) {
	db, _, _ := repository.Get("")

	var myNodeCodes []string
	_ = db.GostNodeBind.Where(db.GostNodeBind.UserCode.Eq(claims.Code)).Pluck(db.GostNodeBind.NodeCode, &myNodeCodes)
	var where = []gen.Condition{
		db.GostNode.Code.In(myNodeCodes...),
	}

	nodes, total, _ := db.GostNode.Where(where...).Order(db.GostNode.IndexValue.Asc(), db.GostNode.Id.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	var dateOnly = time.Now().Format(time.DateOnly)
	for _, node := range nodes {
		var ruleNames []string
		for _, rule := range node.GetRules() {
			getRule := node_rule.Registry.GetRule(rule)
			ruleNames = append(ruleNames, getRule.Name())
		}
		// 获取最近30天的流量
		monthObsInfo := cache2.GetNodeObsDateRange(cache2.MONTH_DATEONLY_LIST, node.Code)
		// 获取今日的流量

		// 获取循环日期的流量
		var obsUseTotal int64
		obsLimit := cache2.GetNodeObsLimit(node.Code)
		if obsLimit.InputBytes == -1 && obsLimit.OutputBytes == -1 {
			obsUseTotal = -1
		} else {
			nowObsInfo := cache2.GetNodeObs(dateOnly, node.Code)
			switch node.LimitKind {
			case model.GOST_NODE_LIMIT_KIND_ALL:
				obsUseTotal += obsLimit.InputBytes + obsLimit.OutputBytes + nowObsInfo.InputBytes + nowObsInfo.OutputBytes
			case model.GOST_NODE_LIMIT_KIND_INPUT:
				obsUseTotal += obsLimit.InputBytes + nowObsInfo.InputBytes
			case model.GOST_NODE_LIMIT_KIND_OUTPUT:
				obsUseTotal += obsLimit.OutputBytes + nowObsInfo.OutputBytes
			}
		}
		list = append(list, Item{
			Code:              node.Code,
			Key:               node.Key,
			Name:              node.Name,
			Remark:            utils.TrinaryOperation(node.Remark == "", "暂无介绍", node.Remark),
			Web:               node.Web,
			Tunnel:            node.Tunnel,
			Forward:           node.Forward,
			Proxy:             node.Proxy,
			P2P:               node.P2P,
			Domain:            node.Domain,
			CustomDomain:      utils.TrinaryOperation(cache2.GetNodeCustomDomain(node.Code), 1, 2),
			DenyDomainPrefix:  node.DenyDomainPrefix,
			UrlTpl:            node.UrlTpl,
			Address:           node.Address,
			ReplaceAddress:    node.ReplaceAddress,
			HttpPort:          node.HttpPort,
			Protocol:          node.Protocol,
			ForwardPorts:      node.ForwardPorts,
			Rules:             node.GetRules(),
			RuleNames:         ruleNames,
			Tags:              node.GetTags(),
			IndexValue:        node.IndexValue,
			P2PDisableForward: node.P2PDisableForward,
			Online:            utils.TrinaryOperation(cache2.GetNodeOnline(node.Code), 1, 2),
			Version:           cache2.GetNodeVersion(node.Code),
			InputBytes:        monthObsInfo.InputBytes,
			OutputBytes:       monthObsInfo.OutputBytes,
			LimitResetIndex:   node.LimitResetIndex,
			LimitTotal:        node.LimitTotal,
			LimitKind:         node.LimitKind,
			LimitUseTotal:     obsUseTotal,
		})
	}
	return list, total
}
