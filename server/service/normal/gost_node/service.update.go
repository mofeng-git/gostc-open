package service

import (
	"errors"
	"go.uber.org/zap"
	"server/model"
	"server/pkg/jwt"
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
	Address               string   `binding:"required" json:"address"`
	Protocol              string   `binding:"required" json:"protocol"`
	Domain                string   `json:"domain"`
	DenyDomainPrefix      string   `json:"denyDomainPrefix"`
	TunnelConnPort        string   `json:"tunnelConnPort"`
	TunnelInPort          string   `json:"tunnelInPort"`
	TunnelMetadata        string   `json:"tunnelMetadata"`
	ForwardConnPort       string   `json:"forwardConnPort"`
	ForwardPorts          string   `json:"forwardPorts"`
	ForwardMetadata       string   `json:"forwardMetadata"`
	TunnelReplaceAddress  string   `json:"tunnelReplaceAddress"`
	ForwardReplaceAddress string   `json:"forwardReplaceAddress"`
	IndexValue            int      `json:"indexValue"`
}

func (service *service) Update(claims jwt.Claims, req UpdateReq) error {
	db, _, log := repository.Get("")
	var nodeBind model.GostNodeBind
	if db.Where("node_code = ? AND user_code = ?", req.Code, claims.Code).First(&nodeBind).RowsAffected == 0 {
		return nil
	}

	var node model.GostNode
	if db.Where("code = ?", req.Code).First(&node).RowsAffected == 0 {
		return errors.New("数据不存在")
	}

	if req.Web == 1 && req.Tunnel != 1 {
		req.Web = 2
	}

	node.Name = req.Name
	node.Remark = req.Remark
	node.Rules = strings.Join(req.Rules, ",")
	node.Tags = strings.Join(req.Tags, ",")

	node.Web = req.Web
	node.Tunnel = req.Tunnel
	node.Forward = req.Forward
	node.Domain = req.Domain
	node.DenyDomainPrefix = req.DenyDomainPrefix
	node.Address = req.Address
	node.Protocol = req.Protocol
	node.TunnelConnPort = req.TunnelConnPort
	node.TunnelInPort = req.TunnelInPort
	node.TunnelMetadata = req.TunnelMetadata
	node.TunnelReplaceAddress = req.TunnelReplaceAddress
	node.ForwardConnPort = req.ForwardConnPort
	node.ForwardPorts = req.ForwardPorts
	node.ForwardMetadata = req.ForwardMetadata
	node.ForwardReplaceAddress = req.ForwardReplaceAddress
	node.IndexValue = req.IndexValue
	if err := db.Save(&node).Error; err != nil {
		log.Error("修改节点失败", zap.Error(err))
		return errors.New("操作失败")
	}
	gost_engine.NodeConfig(db, node.Code)

	var hostCodes []string
	db.Model(&model.GostClientHost{}).Where("node_code = ?", node.Code).Pluck("code", &hostCodes)
	for _, code := range hostCodes {
		gost_engine.ClientHostConfig(db, code)
	}
	var forwardCodes []string
	db.Model(&model.GostClientForward{}).Where("node_code = ?", node.Code).Pluck("code", &forwardCodes)
	for _, code := range forwardCodes {
		gost_engine.ClientForwardConfig(db, code)
	}
	var tunnelCodes []string
	db.Model(&model.GostClientTunnel{}).Where("node_code = ?", node.Code).Pluck("code", &tunnelCodes)
	for _, code := range tunnelCodes {
		gost_engine.ClientTunnelConfig(db, code)
	}
	node_port.Arrange(db, node.Code)
	return nil
}
