package service

import (
	"errors"
	"server/repository"
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

	Domain           string `json:"domain"`
	DenyDomainPrefix string `json:"denyDomainPrefix"`
	Address          string `json:"address"`
	Protocol         string `json:"protocol"`
	TunnelConnPort   string `json:"tunnelConnPort"`
	TunnelInPort     string `json:"tunnelInPort"`
	TunnelMetadata   string `json:"tunnelMetadata"`

	ForwardConnPort string   `json:"forwardConnPort"`
	ForwardPorts    string   `json:"forwardPorts"`
	ForwardMetadata string   `json:"forwardMetadata"`
	Rules           []string `json:"rules"`
	RuleNames       []string `json:"ruleNames"`
	Tags            []string `json:"tags"`

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
		Domain:                node.Domain,
		DenyDomainPrefix:      node.DenyDomainPrefix,
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
		TunnelReplaceAddress:  node.TunnelReplaceAddress,
		ForwardReplaceAddress: node.ForwardReplaceAddress,
		IndexValue:            node.IndexValue,
	}, nil
}
