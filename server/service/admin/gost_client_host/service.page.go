package service

import (
	"gorm.io/gen"
	"server/pkg/bean"
	"server/pkg/utils"
	"server/repository"
	"server/service/common/cache"
	"server/service/common/warn_msg"
	"time"
)

type PageReq struct {
	bean.PageParam
	Account string `json:"account"`
	Enable  int    `json:"enable"`
}

type Item struct {
	UserAccount  string     `json:"userAccount"`
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	TargetIp     string     `json:"targetIp"`
	TargetPort   string     `json:"targetPort"`
	DomainPrefix string     `json:"domainPrefix"`
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
	OnlyChina    int    `json:"onlyChina"`
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

	if req.Enable > 0 {
		where = append(where, db.GostClientHost.Enable.Eq(req.Enable))
	}

	hosts, total, _ := db.GostClientHost.Preload(
		db.GostClientHost.User,
		db.GostClientHost.Client,
		db.GostClientHost.Node,
	).Where(where...).Order(db.GostClientHost.Id.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	for _, host := range hosts {
		obsInfo := cache.GetTunnelObsDateRange(cache.MONTH_DATEONLY_LIST, host.Code)
		list = append(list, Item{
			UserAccount:  host.User.Account,
			Code:         host.Code,
			Name:         host.Name,
			TargetIp:     host.TargetIp,
			TargetPort:   host.TargetPort,
			DomainPrefix: host.DomainPrefix,
			Node: ItemNode{
				Code:    host.NodeCode,
				Name:    host.Node.Name,
				Address: host.Node.Address,
				Online:  utils.TrinaryOperation(cache.GetNodeOnline(host.NodeCode), 1, 2),
				Domain:  host.Node.Domain,
			},
			Client: ItemClient{
				Code:   host.ClientCode,
				Name:   host.Client.Name,
				Online: utils.TrinaryOperation(cache.GetClientOnline(host.ClientCode), 1, 2),
			},
			Config: ItemConfig{
				ChargingType: host.ChargingType,
				Cycle:        host.Cycle,
				Amount:       host.Amount.String(),
				Limiter:      host.Limiter,
				RLimiter:     host.RLimiter,
				CLimiter:     host.CLimiter,
				ExpAt:        time.Unix(host.ExpAt, 0).Format(time.DateTime),
				OnlyChina:    host.OnlyChina,
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
