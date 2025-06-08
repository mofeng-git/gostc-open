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
}

type ItemClient struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Online int    `json:"online"`
}

type ItemNode struct {
	Code              string `json:"code"`
	Name              string `json:"name"`
	Address           string `json:"address"`
	Online            int    `json:"online"`
	P2PDisableForward int    `json:"p2pDisableForward"`
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
		db.GostClientP2P.UserCode.Eq(claims.Code),
	}
	if req.Name != "" {
		where = append(where, db.GostClientP2P.Name.Like("%"+req.Name+"%"))
	}
	if req.Enable > 0 {
		where = append(where, db.GostClientP2P.Enable.Eq(req.Enable))
	}
	if req.ClientCode != "" {
		where = append(where, db.GostClientP2P.ClientCode.Eq(req.ClientCode))
	}
	p2ps, total, _ := db.GostClientP2P.Preload(
		db.GostClientP2P.User,
		db.GostClientP2P.Client,
		db.GostClientP2P.Node,
	).Where(where...).Order(db.GostClientP2P.Id.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	for _, p2p := range p2ps {
		list = append(list, Item{
			Code:       p2p.Code,
			Name:       p2p.Name,
			TargetIp:   p2p.TargetIp,
			TargetPort: p2p.TargetPort,
			VKey:       p2p.VKey,
			Node: ItemNode{
				Code: p2p.NodeCode,
				Name: p2p.Node.Name,
				Address: func() string {
					address, _ := p2p.Node.GetAddress()
					return address
				}(),
				Online:            utils.TrinaryOperation(cache.GetNodeOnline(p2p.NodeCode), 1, 2),
				P2PDisableForward: p2p.Node.P2PDisableForward,
			},
			Client: ItemClient{
				Code:   p2p.ClientCode,
				Name:   p2p.Client.Name,
				Online: utils.TrinaryOperation(cache.GetClientOnline(p2p.ClientCode), 1, 2),
			},
			UserCode:    p2p.UserCode,
			UserAccount: p2p.User.Account,
			Config: ItemConfig{
				ChargingType: p2p.ChargingType,
				Cycle:        p2p.Cycle,
				Amount:       p2p.Amount.String(),
				Limiter:      p2p.Limiter,
				//RLimiter:     p2p.RLimiter,
				//CLimiter:     p2p.CLimiter,
				ExpAt: time.Unix(p2p.ExpAt, 0).Format(time.DateTime),
			},
			Enable:    p2p.Enable,
			WarnMsg:   warn_msg.GetP2PWarnMsg(*p2p),
			CreatedAt: p2p.CreatedAt.Format(time.DateTime),
		})
	}
	return list, total
}
