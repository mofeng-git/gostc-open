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
	Name   string `json:"name"`
	Enable int    `json:"enable"`
}

type Item struct {
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	TargetIp    string     `json:"targetIp"`
	TargetPort  string     `json:"targetPort"`
	VKey        string     `json:"vKey"`
	Node        ItemNode   `json:"node"`
	Client      ItemClient `json:"client"`
	UserCode    string     `json:"userCode"`
	UserAccount string     `json:"userAccount"`
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
	OnlyChina    int    `json:"onlyChina"`
	ExpAt        string `json:"expAt"`
}

func (service *service) Page(claims jwt.Claims, req PageReq) (list []Item, total int64) {
	db, _, _ := repository.Get("")
	var where = []gen.Condition{
		db.GostClientTunnel.UserCode.Eq(claims.Code),
	}
	if req.Name != "" {
		where = append(where, db.GostClientTunnel.Name.Like("%"+req.Name+"%"))
	}
	if req.Enable > 0 {
		where = append(where, db.GostClientTunnel.Enable.Eq(req.Enable))
	}
	tunnels, total, _ := db.GostClientTunnel.Preload(
		db.GostClientTunnel.User,
		db.GostClientTunnel.Client,
		db.GostClientTunnel.Node,
	).Where(where...).Order(db.GostClientTunnel.Id.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	for _, tunnel := range tunnels {
		obsInfo := cache.GetTunnelObsDateRange(cache.MONTH_DATEONLY_LIST, tunnel.Code)
		list = append(list, Item{
			Code:       tunnel.Code,
			Name:       tunnel.Name,
			TargetIp:   tunnel.TargetIp,
			TargetPort: tunnel.TargetPort,
			VKey:       tunnel.VKey,
			Node: ItemNode{
				Code:    tunnel.NodeCode,
				Name:    tunnel.Node.Name,
				Address: tunnel.Node.Address,
				Online:  utils.TrinaryOperation(cache.GetNodeOnline(tunnel.NodeCode), 1, 2),
			},
			Client: ItemClient{
				Code:   tunnel.ClientCode,
				Name:   tunnel.Client.Name,
				Online: utils.TrinaryOperation(cache.GetClientOnline(tunnel.ClientCode), 1, 2),
			},
			UserCode:    tunnel.UserCode,
			UserAccount: tunnel.User.Account,
			Config: ItemConfig{
				ChargingType: tunnel.ChargingType,
				Cycle:        tunnel.Cycle,
				Amount:       tunnel.Amount.String(),
				Limiter:      tunnel.Limiter,
				RLimiter:     tunnel.RLimiter,
				CLimiter:     tunnel.CLimiter,
				ExpAt:        time.Unix(tunnel.ExpAt, 0).Format(time.DateTime),
				OnlyChina:    tunnel.OnlyChina,
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
