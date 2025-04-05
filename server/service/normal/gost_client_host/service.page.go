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
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	TargetIp     string     `json:"targetIp"`
	TargetPort   string     `json:"targetPort"`
	DomainPrefix string     `json:"domainPrefix"`
	DomainFull   string     `json:"domainFull"`
	CustomDomain string     `json:"customDomain"`
	CustomCert   string     `json:"customCert"`
	CustomKey    string     `json:"customKey"`
	CustomEnable int        `json:"customEnable"`
	Node         ItemNode   `json:"node"`
	Client       ItemClient `json:"client"`
	Config       ItemConfig `json:"config"`
	Enable       int        `json:"enable"`
	WarnMsg      string     `json:"warnMsg"`
	CreatedAt    string     `json:"createdAt"`
	InputBytes   int64      `json:"inputBytes"`
	OutputBytes  int64      `json:"outputBytes"`
	WhiteEnable  int        `json:"whiteEnable"`
	BlackEnable  int        `json:"blackEnable"`
	WhiteList    []string   `json:"whiteList"`
	BlackList    []string   `json:"blackList"`
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
	Domain       string `json:"domain"`
	CustomDomain int    `json:"customDomain"`
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
		db.GostClientHost.UserCode.Eq(claims.Code),
	}
	if req.Name != "" {
		where = append(where, db.GostClientHost.Name.Like("%"+req.Name+"%"))
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
			Code:         host.Code,
			Name:         host.Name,
			TargetIp:     host.TargetIp,
			TargetPort:   host.TargetPort,
			DomainPrefix: host.DomainPrefix,
			DomainFull:   host.Node.GetDomainFull(host.DomainPrefix, host.CustomDomain, cache.GetNodeCustomDomain(host.NodeCode)),
			CustomDomain: host.CustomDomain,
			CustomCert:   host.CustomCert,
			CustomKey:    host.CustomKey,
			CustomEnable: utils.TrinaryOperation(cache.GetNodeCustomDomain(host.NodeCode), 1, 2),
			Node: ItemNode{
				Code:         host.NodeCode,
				Name:         host.Node.Name,
				Address:      host.Node.Address,
				Online:       utils.TrinaryOperation(cache.GetNodeOnline(host.NodeCode), 1, 2),
				Domain:       host.Node.Domain,
				CustomDomain: utils.TrinaryOperation(cache.GetNodeCustomDomain(host.NodeCode), 1, 2),
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
			WhiteEnable: host.WhiteEnable,
			BlackEnable: host.BlackEnable,
			WhiteList:   host.GetWhiteList(),
			BlackList:   host.GetBlackList(),
		})
	}
	return list, total
}
