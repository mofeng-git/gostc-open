package service

import (
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"server/model"
	"server/pkg/jwt"
	"server/repository"
	cache2 "server/repository/cache"
	"server/repository/query"
	"server/service/common/node_port"
	"strings"
)

type CreateReq struct {
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

func (service *service) Create(claims jwt.Claims, req CreateReq) error {
	db, _, log := repository.Get("")

	var cfg model.SystemConfigGost
	cache2.GetSystemConfigGost(&cfg)
	if cfg.FuncNode != "1" {
		return errors.New("管理员未启用该功能")
	}

	return db.Transaction(func(tx *query.Query) error {
		var node = model.GostNode{
			Key:                uuid.NewString(),
			Name:               req.Name,
			Remark:             req.Remark,
			Web:                req.Web,
			Tunnel:             req.Tunnel,
			Forward:            req.Forward,
			Proxy:              req.Proxy,
			P2P:                req.P2P,
			Domain:             req.Domain,
			DenyDomainPrefix:   req.DenyDomainPrefix,
			Address:            req.Address,
			ReplaceAddress:     req.ReplaceAddress,
			AllowDomainMatcher: req.AllowDomainMatcher,
			HttpPort:           req.HttpPort,
			UrlTpl:             req.UrlTpl,
			Protocol:           req.Protocol,
			ForwardPorts:       req.ForwardPorts,
			P2PDisableForward:  req.P2PDisableForward,
			Rules:              strings.Join(req.Rules, ","),
			Tags:               strings.Join(req.Tags, ","),
			IndexValue:         req.IndexValue,
			LimitKind:          req.LimitKind,
			LimitTotal:         req.LimitTotal,
			LimitResetIndex:    req.LimitResetIndex,
		}
		if err := tx.GostNode.Create(&node); err != nil {
			log.Error("新增节点失败", zap.Error(err))
			return errors.New("操作失败")
		}
		if err := tx.GostNodeBind.Create(&model.GostNodeBind{
			NodeCode: node.Code,
			UserCode: claims.Code,
		}); err != nil {
			log.Error("新增节点失败", zap.Error(err))
			return errors.New("操作失败")
		}
		if err := tx.GostNodeConfig.Create(&model.GostNodeConfig{
			IndexValue:   1000,
			Name:         "自建节点套餐",
			ChargingType: model.GOST_CONFIG_CHARGING_FREE,
			NodeCode:     node.Code,
		}); err != nil {
			log.Error("新增节点失败", zap.Error(err))
			return errors.New("操作失败")
		}
		node_port.Arrange(tx, node.Code)
		cache2.RefreshNodeObsLimit(node.Code, node.LimitResetIndex)
		cache2.SetNodeInfo(cache2.NodeInfo{
			Code:            node.Code,
			LimitResetIndex: node.LimitResetIndex,
			LimitTotal:      node.LimitTotal,
			LimitKind:       node.LimitKind,
		})
		return nil
	})
}
