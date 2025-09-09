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
	Name         string `json:"name"`
	Account      string `json:"account"`
	ClientName   string `json:"clientName"`
	NodeName     string `json:"nodeName"`
	Enable       int    `json:"enable"`
	ClientOnline int    `json:"clientOnline"`
	NodeOnline   int    `json:"nodeOnline"`
}

type Item struct {
	UserAccount  string     `json:"userAccount"`
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	TargetIp     string     `json:"targetIp"`
	TargetPort   string     `json:"targetPort"`
	DomainPrefix string     `json:"domainPrefix"`
	DomainFull   string     `json:"domainFull"`
	Node         ItemNode   `json:"node"`
	Client       ItemClient `json:"client"`
	Config       ItemConfig `json:"config"`
	Enable       int        `json:"enable"`
	WarnMsg      string     `json:"warnMsg"`
	CreatedAt    string     `json:"createdAt"`
	InputBytes   int64      `json:"inputBytes"`
	OutputBytes  int64      `json:"outputBytes"`
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
	Domain  string `json:"domain"`
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
		where = append(where, db.GostClientHost.UserCode.In(userCodes...))
	}
	if req.NodeName != "" {
		var codes []string
		_ = db.GostNode.Where(db.GostNode.Name.Like("%"+req.NodeName+"%")).Pluck(db.GostNode.Code, &codes)
		where = append(where, db.GostClientHost.NodeCode.In(codes...))
	}
	if req.Enable > 0 {
		where = append(where, db.GostClientHost.Enable.Eq(req.Enable))
	}
	if req.Name != "" {
		where = append(where, db.GostClientHost.Name.Like("%"+req.Name+"%"))
	}
	if req.ClientName != "" {
		var clientCodes []string
		_ = db.GostClient.Where(db.GostClient.Name.Like("%"+req.ClientName+"%")).Pluck(db.GostClient.Code, &clientCodes)
		where = append(where, db.GostClientHost.ClientCode.In(clientCodes...))
	}
	if req.ClientOnline > 0 {
		var clientCodes []string
		_ = db.GostClient.Pluck(db.GostClient.Code, &clientCodes)
		var onlineCodes []string
		var offlineCodes []string
		for _, code := range clientCodes {
			if cache2.GetClientOnline(code) {
				onlineCodes = append(onlineCodes, code)
			} else {
				offlineCodes = append(offlineCodes, code)
			}
		}
		if req.ClientOnline == 1 {
			where = append(where, db.GostClientHost.ClientCode.In(onlineCodes...))
		} else {
			where = append(where, db.GostClientHost.ClientCode.In(offlineCodes...))
		}
	}

	if req.NodeOnline > 0 {
		var nodeCodes []string
		_ = db.GostNode.Pluck(db.GostNode.Code, &nodeCodes)
		var onlineCodes []string
		var offlineCodes []string
		for _, code := range nodeCodes {
			if cache2.GetNodeOnline(code) {
				onlineCodes = append(onlineCodes, code)
			} else {
				offlineCodes = append(offlineCodes, code)
			}
		}
		if req.NodeOnline == 1 {
			where = append(where, db.GostClientHost.NodeCode.In(onlineCodes...))
		} else {
			where = append(where, db.GostClientHost.NodeCode.In(offlineCodes...))
		}
	}

	hosts, total, _ := db.GostClientHost.Preload(
		db.GostClientHost.User,
		db.GostClientHost.Client,
		db.GostClientHost.Node,
	).Where(where...).Order(db.GostClientHost.Id.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	for _, host := range hosts {
		obsInfo := cache2.GetTunnelObsDateRange(cache2.MONTH_DATEONLY_LIST, host.Code)
		list = append(list, Item{
			UserAccount:  host.User.Account,
			Code:         host.Code,
			Name:         host.Name,
			TargetIp:     host.TargetIp,
			TargetPort:   host.TargetPort,
			DomainPrefix: host.DomainPrefix,
			DomainFull:   host.Node.GetDomainFull(host.DomainPrefix, host.GetCustomDomain(), cache2.GetNodeCustomDomain(host.NodeCode)),
			Node: ItemNode{
				Code: host.NodeCode,
				Name: host.Node.Name,
				Address: func() string {
					address, _, _ := net.SplitHostPort(host.Node.Address)
					return address
				}(),
				Online: utils.TrinaryOperation(cache2.GetNodeOnline(host.NodeCode), 1, 2),
				Domain: host.Node.Domain,
			},
			Client: ItemClient{
				Code:   host.ClientCode,
				Name:   host.Client.Name,
				Online: utils.TrinaryOperation(cache2.GetClientOnline(host.ClientCode), 1, 2),
			},
			Config: ItemConfig{
				ChargingType: host.ChargingType,
				Cycle:        host.Cycle,
				Amount:       host.Amount.String(),
				Limiter:      host.Limiter,
				//RLimiter:     host.RLimiter,
				//CLimiter:     host.CLimiter,
				ExpAt: time.Unix(host.ExpAt, 0).Format(time.DateTime),
			},
			Enable:      host.Enable,
			WarnMsg:     warn_msg.GetHostWarnMsg(*host),
			CreatedAt:   host.CreatedAt.Format(time.DateTime),
			InputBytes:  obsInfo.InputBytes,
			OutputBytes: obsInfo.OutputBytes,
		})
	}
	return list, total
}
