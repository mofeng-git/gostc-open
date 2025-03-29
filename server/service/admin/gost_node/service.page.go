package service

import (
	"gorm.io/gen"
	"server/pkg/bean"
	"server/pkg/utils"
	"server/repository"
	"server/service/common/cache"
	"server/service/common/node_rule"
)

type PageReq struct {
	bean.PageParam
	Name string `json:"name"`
	Bind int    `json:"bind"`
}

type Item struct {
	Code        string `json:"code"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Remark      string `json:"remark"`
	UserAccount string `json:"userAccount"`

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

func (service *service) Page(req PageReq) (list []Item, total int64) {
	db, _, _ := repository.Get("")
	var where []gen.Condition
	if req.Name != "" {
		where = append(where, db.GostNode.Name.Like("%"+req.Name+"%"))
	}

	var nodeCodes []string
	_ = db.GostNodeBind.Pluck(db.GostNodeBind.NodeCode, &nodeCodes)
	switch req.Bind {
	case 1:
		// 用户节点
		where = append(where, db.GostNode.Code.In(nodeCodes...))
	case 2:
		// 系统节点
		where = append(where, db.GostNode.Code.NotIn(nodeCodes...))
	}
	nodes, total, _ := db.GostNode.Where(where...).Order(db.GostNode.IndexValue.Asc(), db.GostNode.Id.Desc()).FindByPage(req.GetOffset(), req.GetLimit())

	// 查询节点的绑定用户
	nodeBinds, _ := db.GostNodeBind.Preload(db.GostNodeBind.User).Find()
	var nodeBindAccountMap = make(map[string]string)
	for _, bind := range nodeBinds {
		nodeBindAccountMap[bind.NodeCode] = bind.User.Account
	}
	for _, node := range nodes {
		var ruleNames []string
		for _, rule := range node.GetRules() {
			ruleNames = append(ruleNames, node_rule.RuleMap[rule].Name())
		}
		name := node.Name
		account := nodeBindAccountMap[node.Code]
		obsInfo := cache.GetNodeObsDateRange(cache.MONTH_DATEONLY_LIST, node.Code)
		list = append(list, Item{
			Code:                  node.Code,
			Key:                   node.Key,
			Name:                  name,
			Remark:                utils.TrinaryOperation(node.Remark == "", "暂无介绍", node.Remark),
			UserAccount:           account,
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
