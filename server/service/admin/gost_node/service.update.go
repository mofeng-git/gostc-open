package service

import (
	"errors"
	"go.uber.org/zap"
	"server/repository"
	"server/service/common/node_port"
	"server/service/gost_engine"
	"strings"
)

type UpdateReq struct {
	Code                  string   `binding:"required" json:"code" label:"编号"`
	Name                  string   `binding:"required" json:"name" label:"名称"`
	Remark                string   `json:"remark"`
	Rules                 []string `json:"rules"`
	Tags                  []string `json:"tags"`
	Web                   int      `binding:"required" json:"web" label:"是否启用域名解析"`
	Tunnel                int      `binding:"required" json:"tunnel" label:"是否启用私有隧道"`
	Forward               int      `binding:"required" json:"forward" label:"是否启用端口转发"`
	Proxy                 int      `binding:"required" json:"proxy" label:"是否启用代理隧道"`
	P2P                   int      `binding:"required" json:"p2p" label:"是否启用P2P隧道"`
	Address               string   `binding:"required" json:"address"`
	Protocol              string   `binding:"required" json:"protocol"`
	Domain                string   `json:"domain"`
	DenyDomainPrefix      string   `json:"denyDomainPrefix"`
	UrlTpl                string   `json:"urlTpl"`
	TunnelConnPort        string   `json:"tunnelConnPort"`
	TunnelInPort          string   `json:"tunnelInPort"`
	TunnelMetadata        string   `json:"tunnelMetadata"`
	ForwardConnPort       string   `json:"forwardConnPort"`
	ForwardPorts          string   `json:"forwardPorts"`
	ForwardMetadata       string   `json:"forwardMetadata"`
	TunnelReplaceAddress  string   `json:"tunnelReplaceAddress"`
	ForwardReplaceAddress string   `json:"forwardReplaceAddress"`
	P2PPort               string   `json:"p2pPort"`
	IndexValue            int      `json:"indexValue"`
}

func (service *service) Update(req UpdateReq) error {
	db, _, log := repository.Get("")
	node, err := db.GostNode.Where(db.GostNode.Code.Eq(req.Code)).First()
	if err != nil {
		return errors.New("节点不存在")
	}
	if req.Web == 1 && req.Tunnel != 1 {
		req.Web = 2
	}
	if req.Proxy == 1 && req.Forward != 1 {
		req.Proxy = 2
	}

	node.Name = req.Name
	node.Remark = req.Remark
	node.Rules = strings.Join(req.Rules, ",")
	node.Tags = strings.Join(req.Tags, ",")

	node.Web = req.Web
	node.Tunnel = req.Tunnel
	node.Forward = req.Forward
	node.Proxy = req.Proxy
	node.P2P = req.P2P
	node.Domain = req.Domain
	node.DenyDomainPrefix = req.DenyDomainPrefix
	node.Address = req.Address
	node.Protocol = req.Protocol
	node.UrlTpl = req.UrlTpl
	node.TunnelConnPort = req.TunnelConnPort
	node.TunnelInPort = req.TunnelInPort
	node.TunnelMetadata = req.TunnelMetadata
	node.TunnelReplaceAddress = req.TunnelReplaceAddress
	node.ForwardConnPort = req.ForwardConnPort
	node.ForwardPorts = req.ForwardPorts
	node.ForwardMetadata = req.ForwardMetadata
	node.ForwardReplaceAddress = req.ForwardReplaceAddress
	node.P2PPort = req.P2PPort
	node.IndexValue = req.IndexValue
	if err := db.GostNode.Save(node); err != nil {
		log.Error("修改节点失败", zap.Error(err))
		return errors.New("操作失败")
	}
	gost_engine.NodeConfig(db, node.Code)

	var hostCodes []string
	_ = db.GostClientHost.Where(db.GostClientHost.NodeCode.Eq(node.Code)).Pluck(db.GostClientHost.Code, &hostCodes)
	for _, code := range hostCodes {
		gost_engine.ClientHostConfig(db, code)
	}
	var forwardCodes []string
	_ = db.GostClientForward.Where(db.GostClientForward.NodeCode.Eq(node.Code)).Pluck(db.GostClientForward.Code, &forwardCodes)
	for _, code := range forwardCodes {
		gost_engine.ClientForwardConfig(db, code)
	}
	var tunnelCodes []string
	_ = db.GostClientTunnel.Where(db.GostClientTunnel.NodeCode.Eq(node.Code)).Pluck(db.GostClientTunnel.Code, &tunnelCodes)
	for _, code := range tunnelCodes {
		gost_engine.ClientTunnelConfig(db, code)
	}
	var proxyCodes []string
	_ = db.GostClientProxy.Where(db.GostClientProxy.NodeCode.Eq(node.Code)).Pluck(db.GostClientProxy.Code, &proxyCodes)
	for _, code := range proxyCodes {
		gost_engine.ClientProxyConfig(db, code)
	}
	var p2pCodes []string
	_ = db.GostClientP2P.Where(db.GostClientP2P.NodeCode.Eq(node.Code)).Pluck(db.GostClientP2P.Code, &p2pCodes)
	for _, code := range p2pCodes {
		gost_engine.ClientP2PConfig(db, code)
	}
	node_port.Arrange(db, node.Code)
	return nil
}
