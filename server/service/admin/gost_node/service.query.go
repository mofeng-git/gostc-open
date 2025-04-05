package service

import (
	"errors"
	"server/pkg/utils"
	"server/repository"
	"server/service/common/cache"
	"server/service/common/node_rule"
)

type QueryReq struct {
	Code string `binding:"required" json:"code" label:"编号"`
}

type QueryResp struct {
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
	Address          string `json:"address"`
	Protocol         string `json:"protocol"`
	TunnelConnPort   string `json:"tunnelConnPort"`
	TunnelInPort     string `json:"tunnelInPort"`
	TunnelMetadata   string `json:"tunnelMetadata"`

	ForwardConnPort string `json:"forwardConnPort"`
	ForwardPorts    string `json:"forwardPorts"`
	ForwardMetadata string `json:"forwardMetadata"`

	P2PPort string `json:"p2pPort"`

	Rules     []string `json:"rules"`
	RuleNames []string `json:"ruleNames"`
	Tags      []string `json:"tags"`

	TunnelReplaceAddress  string `json:"tunnelReplaceAddress"`
	ForwardReplaceAddress string `json:"forwardReplaceAddress"`
	IndexValue            int    `json:"indexValue"`
}

func (service *service) Query(req QueryReq) (QueryResp, error) {
	db, _, _ := repository.Get("")
	node, err := db.GostNode.Where(db.GostNode.Code.Eq(req.Code)).First()
	if err != nil {
		return QueryResp{}, errors.New("节点不存在")
	}
	var ruleNames []string
	for _, rule := range node.GetRules() {
		ruleNames = append(ruleNames, node_rule.RuleMap[rule].Name())
	}
	return QueryResp{
		Code:                  node.Code,
		Key:                   node.Key,
		Name:                  node.Name,
		Remark:                node.Remark,
		Web:                   node.Web,
		Tunnel:                node.Tunnel,
		Forward:               node.Forward,
		Proxy:                 node.Proxy,
		P2P:                   node.P2P,
		Domain:                node.Domain,
		CustomDomain:          utils.TrinaryOperation(cache.GetNodeCustomDomain(node.Code), 1, 2),
		DenyDomainPrefix:      node.DenyDomainPrefix,
		Address:               node.Address,
		Protocol:              node.Protocol,
		TunnelConnPort:        node.TunnelConnPort,
		TunnelInPort:          node.TunnelInPort,
		TunnelMetadata:        node.TunnelMetadata,
		ForwardConnPort:       node.ForwardConnPort,
		ForwardPorts:          node.ForwardPorts,
		ForwardMetadata:       node.ForwardMetadata,
		P2PPort:               node.P2PPort,
		Rules:                 node.GetRules(),
		RuleNames:             ruleNames,
		Tags:                  node.GetTags(),
		TunnelReplaceAddress:  node.TunnelReplaceAddress,
		ForwardReplaceAddress: node.ForwardReplaceAddress,
		IndexValue:            node.IndexValue,
	}, nil
}
