package service

import (
	"gorm.io/gen"
	"net"
	"server/pkg/bean"
	"server/pkg/utils"
	"server/repository"
	cache2 "server/repository/cache"
	"server/service/common/warn_msg"
	"time"
)

type PageReq struct {
	bean.PageParam
	Name       string `json:"name"`
	Account    string `json:"account"`
	ClientName string `json:"clientName"`
	NodeName   string `json:"nodeName"`
	Enable     int    `json:"enable"`
}

type Item struct {
	UserAccount string     `json:"userAccount"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	TargetIp    string     `json:"targetIp"`
	TargetPort  string     `json:"targetPort"`
	VKey        string     `json:"vKey"`
	Node        ItemNode   `json:"node"`
	Client      ItemClient `json:"client"`
	Config      ItemConfig `json:"config"`
	Enable      int        `json:"enable"`
	WarnMsg     string     `json:"warnMsg"`
	CreatedAt   string     `json:"createdAt"`
	InputBytes  int64      `json:"inputBytes"`
	OutputBytes int64      `json:"outputBytes"`
}

type ItemClient struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Online int    `json:"online"`
}

type ItemNode struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Online  int    `json:"online"`
}

type ItemConfig struct {
	ChargingType int    `json:"chargingType"`
	Cycle        int    `json:"cycle"`
	Amount       string `json:"amount"`
	Limiter      int    `json:"limiter"`
	RLimiter     int    `json:"rLimiter"`
	CLimiter     int    `json:"cLimiter"`
	ExpAt        string `json:"expAt"`
}

func (service *service) Page(req PageReq) (list []Item, total int64) {
	db, _, _ := repository.Get("")
	var where []gen.Condition
	if req.Account != "" {
		var userCodes []string
		_ = db.SystemUser.Where(db.SystemUser.Account.Like("%"+req.Account+"%")).Pluck(db.SystemUser.Code, &userCodes)
		where = append(where, db.GostClientTunnel.UserCode.In(userCodes...))
	}
	if req.NodeName != "" {
		var codes []string
		_ = db.GostNode.Where(db.GostNode.Name.Like("%"+req.NodeName+"%")).Pluck(db.GostNode.Code, &codes)
		where = append(where, db.GostClientTunnel.NodeCode.In(codes...))
	}
	if req.Enable > 0 {
		where = append(where, db.GostClientTunnel.Enable.Eq(req.Enable))
	}
	if req.Name != "" {
		where = append(where, db.GostClientTunnel.Name.Like("%"+req.Name+"%"))
	}
	if req.ClientName != "" {
		var clientCodes []string
		_ = db.GostClient.Where(db.GostClient.Name.Like("%"+req.ClientName+"%")).Pluck(db.GostClient.Code, &clientCodes)
		where = append(where, db.GostClientTunnel.ClientCode.In(clientCodes...))
	}
	tunnels, total, _ := db.GostClientTunnel.Preload(
		db.GostClientTunnel.User,
		db.GostClientTunnel.Client,
		db.GostClientTunnel.Node,
	).Where(where...).Order(db.GostClientTunnel.Id.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	for _, tunnel := range tunnels {
		obsInfo := cache2.GetTunnelObsDateRange(cache2.MONTH_DATEONLY_LIST, tunnel.Code)
		list = append(list, Item{
			Code:       tunnel.Code,
			Name:       tunnel.Name,
			TargetIp:   tunnel.TargetIp,
			TargetPort: tunnel.TargetPort,
			VKey:       tunnel.VKey,
			Node: ItemNode{
				Code: tunnel.NodeCode,
				Name: tunnel.Node.Name,
				Address: func() string {
					address, _, _ := net.SplitHostPort(tunnel.Node.Address)
					return address
				}(),
				Online: utils.TrinaryOperation(cache2.GetNodeOnline(tunnel.NodeCode), 1, 2),
			},
			Client: ItemClient{
				Code:   tunnel.ClientCode,
				Name:   tunnel.Client.Name,
				Online: utils.TrinaryOperation(cache2.GetClientOnline(tunnel.ClientCode), 1, 2),
			},
			UserAccount: tunnel.User.Account,
			Config: ItemConfig{
				ChargingType: tunnel.ChargingType,
				Cycle:        tunnel.Cycle,
				Amount:       tunnel.Amount.String(),
				Limiter:      tunnel.Limiter,
				//RLimiter:     tunnel.RLimiter,
				//CLimiter:     tunnel.CLimiter,
				ExpAt: time.Unix(tunnel.ExpAt, 0).Format(time.DateTime),
			},
			Enable:      tunnel.Enable,
			WarnMsg:     warn_msg.GetTunnelWarnMsg(*tunnel),
			CreatedAt:   tunnel.CreatedAt.Format(time.DateTime),
			InputBytes:  obsInfo.InputBytes,
			OutputBytes: obsInfo.OutputBytes,
		})
	}
	return list, total
}
