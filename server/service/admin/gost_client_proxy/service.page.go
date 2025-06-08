package service

import (
	"gorm.io/gen"
	"net"
	"server/pkg/bean"
	"server/pkg/utils"
	"server/repository"
	"server/service/common/cache"
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
	Port        string     `json:"port"`
	Protocol    string     `json:"protocol"`
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
		where = append(where, db.GostClientProxy.UserCode.In(userCodes...))
	}
	if req.NodeName != "" {
		var codes []string
		_ = db.GostNode.Where(db.GostNode.Name.Like("%"+req.NodeName+"%")).Pluck(db.GostNode.Code, &codes)
		where = append(where, db.GostClientProxy.NodeCode.In(codes...))
	}
	if req.Enable > 0 {
		where = append(where, db.GostClientProxy.Enable.Eq(req.Enable))
	}
	if req.Name != "" {
		where = append(where, db.GostClientProxy.Name.Like("%"+req.Name+"%"))
	}
	if req.ClientName != "" {
		var clientCodes []string
		_ = db.GostClient.Where(db.GostClient.Name.Like("%"+req.ClientName+"%")).Pluck(db.GostClient.Code, &clientCodes)
		where = append(where, db.GostClientProxy.ClientCode.In(clientCodes...))
	}
	proxys, total, _ := db.GostClientProxy.Preload(
		db.GostClientProxy.User,
		db.GostClientProxy.Client,
		db.GostClientProxy.Node,
	).Where(where...).Order(db.GostClientProxy.Id.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	for _, proxy := range proxys {
		obsInfo := cache.GetTunnelObsDateRange(cache.MONTH_DATEONLY_LIST, proxy.Code)
		list = append(list, Item{
			UserAccount: proxy.User.Account,
			Code:        proxy.Code,
			Name:        proxy.Name,
			Port:        proxy.Port,
			Protocol:    proxy.Protocol,
			Node: ItemNode{
				Code: proxy.NodeCode,
				Name: proxy.Node.Name,
				Address: func() string {
					address, _, _ := net.SplitHostPort(proxy.Node.Address)
					return address
				}(),
				Online: utils.TrinaryOperation(cache.GetNodeOnline(proxy.NodeCode), 1, 2),
			},
			Client: ItemClient{
				Code:   proxy.ClientCode,
				Name:   proxy.Client.Name,
				Online: utils.TrinaryOperation(cache.GetClientOnline(proxy.ClientCode), 1, 2),
			},
			Config: ItemConfig{
				ChargingType: proxy.ChargingType,
				Cycle:        proxy.Cycle,
				Amount:       proxy.Amount.String(),
				Limiter:      proxy.Limiter,
				//RLimiter:     proxy.RLimiter,
				//CLimiter:     proxy.CLimiter,
				ExpAt: time.Unix(proxy.ExpAt, 0).Format(time.DateTime),
			},
			Enable:      proxy.Enable,
			WarnMsg:     warn_msg.GetProxyWarnMsg(*proxy),
			CreatedAt:   proxy.CreatedAt.Format(time.DateTime),
			InputBytes:  obsInfo.InputBytes,
			OutputBytes: obsInfo.OutputBytes,
		})
	}
	return list, total
}
