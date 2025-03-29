package service

import (
	"gorm.io/gen"
	"server/pkg/bean"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/service/common/cache"
	"server/service/common/node_rule"
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
	DenyDomainPrefix string `json:"denyDomainPrefix"`
	UrlTpl           string `json:"urlTpl"`
	Address          string `json:"address"`
	Protocol         string `json:"protocol"`
	TunnelConnPort   string `json:"tunnelConnPort"`
	TunnelInPort     string `json:"tunnelInPort"`
	TunnelMetadata   string `json:"tunnelMetadata"`

	ForwardConnPort       string   `json:"forwardConnPort"`
	ForwardPorts          string   `json:"forwardPorts"`
	ForwardMetadata       string   `json:"forwardMetadata"`
	Rules                 []string `json:"rules"`
	RuleNames             []string `json:"ruleNames"`
	Tags                  []string `json:"tags"`
	IndexValue            int      `json:"indexValue"`
	TunnelReplaceAddress  string   `json:"tunnelReplaceAddress"`
	ForwardReplaceAddress string   `json:"forwardReplaceAddress"`
	P2PPort               string   `json:"p2pPort"`

	Online      int    `json:"online"`
	Version     string `json:"version"`
	InputBytes  int64  `json:"inputBytes"`
	OutputBytes int64  `json:"outputBytes"`
}

func (service *service) Page(claims jwt.Claims, req PageReq) (list []Item, total int64) {
	db, _, _ := repository.Get("")

	var myNodeCodes []string
	_ = db.GostNodeBind.Where(db.GostNodeBind.UserCode.Eq(claims.Code)).Pluck(db.GostNodeBind.NodeCode, &myNodeCodes)
	var where = []gen.Condition{
		db.GostNode.Code.In(myNodeCodes...),
	}

	nodes, total, _ := db.GostNode.Where(where...).Order(db.GostNode.IndexValue.Asc(), db.GostNode.Id.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	for _, node := range nodes {
		var ruleNames []string
		for _, rule := range node.GetRules() {
			ruleNames = append(ruleNames, node_rule.RuleMap[rule].Name())
		}
		obsInfo := cache.GetNodeObsDateRange(cache.MONTH_DATEONLY_LIST, node.Code)
		list = append(list, Item{
			Code:                  node.Code,
			Key:                   node.Key,
			Name:                  node.Name,
			Remark:                utils.TrinaryOperation(node.Remark == "", "暂无介绍", node.Remark),
			Web:                   node.Web,
			Tunnel:                node.Tunnel,
			Forward:               node.Forward,
			Proxy:                 node.Proxy,
			P2P:                   node.P2P,
			Domain:                node.Domain,
			DenyDomainPrefix:      node.DenyDomainPrefix,
			UrlTpl:                node.UrlTpl,
			Address:               node.Address,
			Protocol:              node.Protocol,
			TunnelConnPort:        node.TunnelConnPort,
			TunnelInPort:          node.TunnelInPort,
			TunnelMetadata:        node.TunnelMetadata,
			ForwardConnPort:       node.ForwardConnPort,
			ForwardPorts:          node.ForwardPorts,
			ForwardMetadata:       node.ForwardMetadata,
			Rules:                 node.GetRules(),
			RuleNames:             ruleNames,
			Tags:                  node.GetTags(),
			IndexValue:            node.IndexValue,
			TunnelReplaceAddress:  node.TunnelReplaceAddress,
			ForwardReplaceAddress: node.ForwardReplaceAddress,
			P2PPort:               node.P2PPort,
			Online:                utils.TrinaryOperation(cache.GetNodeOnline(node.Code), 1, 2),
			Version:               cache.GetNodeVersion(node.Code),
			InputBytes:            obsInfo.InputBytes,
			OutputBytes:           obsInfo.OutputBytes,
		})
	}
	return list, total
}
