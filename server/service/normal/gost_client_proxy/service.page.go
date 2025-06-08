package service

import (
	"gorm.io/gen"
	"server/pkg/bean"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/service/common/cache"
	"server/service/common/warn_msg"
	"time"
)

type PageReq struct {
	bean.PageParam
	Name       string `json:"name"`
	Enable     int    `json:"enable"`
	ClientCode string `json:"clientCode"`
}

type Item struct {
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Port        string     `json:"port"`
	Protocol    string     `json:"protocol"`
	AuthUser    string     `json:"authUser"`
	AuthPwd     string     `json:"authPwd"`
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
	Code         string `json:"code"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	Online       int    `json:"online"`
	ForwardPorts string `json:"forwardPorts"`
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

func (service *service) Page(claims jwt.Claims, req PageReq) (list []Item, total int64) {
	db, _, _ := repository.Get("")
	var where = []gen.Condition{
		db.GostClientProxy.UserCode.Eq(claims.Code),
	}
	if req.Name != "" {
		where = append(where, db.GostClientProxy.Name.Like("%"+req.Name+"%"))
	}
	if req.Enable > 0 {
		where = append(where, db.GostClientProxy.Enable.Eq(req.Enable))
	}
	if req.ClientCode != "" {
		where = append(where, db.GostClientProxy.ClientCode.Eq(req.ClientCode))
	}
	proxys, total, _ := db.GostClientProxy.Preload(
		db.GostClientProxy.User,
		db.GostClientProxy.Client,
		db.GostClientProxy.Node,
	).Where(where...).Order(db.GostClientProxy.Id.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	for _, proxy := range proxys {
		obsInfo := cache.GetTunnelObsDateRange(cache.MONTH_DATEONLY_LIST, proxy.Code)
		list = append(list, Item{
			Code:     proxy.Code,
			Name:     proxy.Name,
			Port:     proxy.Port,
			AuthUser: proxy.AuthUser,
			AuthPwd:  proxy.AuthPwd,
			Protocol: proxy.Protocol,
			Node: ItemNode{
				Code: proxy.NodeCode,
				Name: proxy.Node.Name,
				Address: func() string {
					address, _ := proxy.Node.GetAddress()
					return address
				}(),
				Online:       utils.TrinaryOperation(cache.GetNodeOnline(proxy.NodeCode), 1, 2),
				ForwardPorts: proxy.Node.ForwardPorts,
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
