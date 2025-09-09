package service

import (
	"errors"
	"go.uber.org/zap"
	"server/pkg/jwt"
	"server/repository"
	cache2 "server/repository/cache"
	"server/service/common/node_port"
	"server/service/engine"
	"strings"
)

type UpdateReq struct {
	Code               string   `binding:"required" json:"code" label:"编号"`
	Name               string   `binding:"required" json:"name" label:"名称"`
	Remark             string   `json:"remark"`
	Rules              []string `json:"rules"`
	Tags               []string `json:"tags"`
	Web                int      `binding:"required" json:"web" label:"是否启用域名解析"`
	Tunnel             int      `binding:"required" json:"tunnel" label:"是否启用私有隧道"`
	Forward            int      `binding:"required" json:"forward" label:"是否启用端口转发"`
	Proxy              int      `binding:"required" json:"proxy" label:"是否启用代理隧道"`
	P2P                int      `binding:"required" json:"p2p" label:"是否启用P2P隧道"`
	Address            string   `binding:"required" json:"address"`
	ReplaceAddress     string   `json:"replaceAddress"`
	Protocol           string   `binding:"required" json:"protocol"`
	Domain             string   `json:"domain"`
	DenyDomainPrefix   string   `json:"denyDomainPrefix"`
	AllowDomainMatcher int      `json:"allowDomainMatcher"`
	UrlTpl             string   `json:"urlTpl"`
	HttpPort           string   `json:"httpPort"`
	ForwardPorts       string   `json:"forwardPorts"`
	P2PDisableForward  int      `json:"p2pDisableForward"`
	IndexValue         int      `json:"indexValue"`

	LimitResetIndex int `json:"limitResetIndex"`
	LimitTotal      int `json:"limitTotal"`
	LimitKind       int `json:"limitKind"`
}

func (service *service) Update(claims jwt.Claims, req UpdateReq) error {
	db, _, log := repository.Get("")
	nodeBind, _ := db.GostNodeBind.Where(db.GostNodeBind.NodeCode.Eq(req.Code), db.GostNodeBind.UserCode.Eq(claims.Code)).First()
	if nodeBind == nil {
		return nil
	}

	node, _ := db.GostNode.Where(db.GostNode.Code.Eq(req.Code)).First()
	if node == nil {
		return errors.New("数据不存在")
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
	node.ReplaceAddress = req.ReplaceAddress
	node.AllowDomainMatcher = req.AllowDomainMatcher
	node.HttpPort = req.HttpPort
	node.Protocol = req.Protocol
	node.UrlTpl = req.UrlTpl
	node.ForwardPorts = req.ForwardPorts
	node.P2PDisableForward = req.P2PDisableForward
	node.IndexValue = req.IndexValue

	node.LimitResetIndex = req.LimitResetIndex
	node.LimitTotal = req.LimitTotal
	node.LimitKind = req.LimitKind

	if err := db.GostNode.Save(node); err != nil {
		log.Error("修改节点失败", zap.Error(err))
		return errors.New("操作失败")
	}
	engine.NodeConfig(db, node.Code)
	engine.ClientAllConfigUpdateByNodeCode(db, node.Code)
	node_port.Arrange(db, node.Code)
	cache2.RefreshNodeObsLimit(node.Code, node.LimitResetIndex)
	cache2.SetNodeInfo(cache2.NodeInfo{
		Code:            node.Code,
		LimitResetIndex: node.LimitResetIndex,
		LimitTotal:      node.LimitTotal,
		LimitKind:       node.LimitKind,
	})
	return nil
}
