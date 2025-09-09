package service

import (
	"gorm.io/gen"
	"server/pkg/bean"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	cache2 "server/repository/cache"
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
	Code                string     `json:"code"`
	Name                string     `json:"name"`
	TargetIp            string     `json:"targetIp"`
	TargetPort          string     `json:"targetPort"`
	TargetHttps         int        `json:"targetHttps"`
	DomainPrefix        string     `json:"domainPrefix"`
	DomainFull          string     `json:"domainFull"`
	CustomDomain        string     `json:"customDomain"`
	CustomCert          string     `json:"customCert"`
	CustomKey           string     `json:"customKey"`
	CustomEnable        int        `json:"customEnable"`
	CustomForceHttps    int        `json:"customForceHttps"`
	CustomDomainMatcher int        `json:"customDomainMatcher"`
	Node                ItemNode   `json:"node"`
	Client              ItemClient `json:"client"`
	Config              ItemConfig `json:"config"`
	Enable              int        `json:"enable"`
	WarnMsg             string     `json:"warnMsg"`
	CreatedAt           string     `json:"createdAt"`
	InputBytes          int64      `json:"inputBytes"`
	OutputBytes         int64      `json:"outputBytes"`
	WhiteEnable         int        `json:"whiteEnable"`
	WhiteList           []string   `json:"whiteList"`
}

type ItemClient struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Online int    `json:"online"`
}

type ItemNode struct {
	Code               string `json:"code"`
	Name               string `json:"name"`
	Address            string `json:"address"`
	Online             int    `json:"online"`
	Domain             string `json:"domain"`
	CustomDomain       int    `json:"customDomain"`
	AllowDomainMatcher int    `json:"allowDomainMatcher"`
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
		db.GostClientHost.UserCode.Eq(claims.Code),
	}
	if req.Name != "" {
		where = append(where, db.GostClientHost.Name.Like("%"+req.Name+"%"))
	}
	if req.Enable > 0 {
		where = append(where, db.GostClientHost.Enable.Eq(req.Enable))
	}
	if req.ClientCode != "" {
		where = append(where, db.GostClientHost.ClientCode.Eq(req.ClientCode))
	}
	hosts, total, _ := db.GostClientHost.Preload(
		db.GostClientHost.User,
		db.GostClientHost.Client,
		db.GostClientHost.Node,
	).Where(where...).Order(db.GostClientHost.Id.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	for _, host := range hosts {
		obsInfo := cache2.GetTunnelObsDateRange(cache2.MONTH_DATEONLY_LIST, host.Code)
		list = append(list, Item{
			Code:                host.Code,
			Name:                host.Name,
			TargetIp:            host.TargetIp,
			TargetPort:          host.TargetPort,
			TargetHttps:         host.TargetHttps,
			DomainPrefix:        host.DomainPrefix,
			DomainFull:          host.Node.GetDomainFull(host.DomainPrefix, host.GetCustomDomain(), cache2.GetNodeCustomDomain(host.NodeCode)),
			CustomDomain:        host.CustomDomain,
			CustomCert:          host.CustomCert,
			CustomKey:           host.CustomKey,
			CustomForceHttps:    host.CustomForceHttps,
			CustomDomainMatcher: host.CustomDomainMatcher,
			CustomEnable:        utils.TrinaryOperation(cache2.GetNodeCustomDomain(host.NodeCode), 1, 2),
			Node: ItemNode{
				Code: host.NodeCode,
				Name: host.Node.Name,
				Address: func() string {
					address, _ := host.Node.GetAddress()
					return address
				}(),
				Online:             utils.TrinaryOperation(cache2.GetNodeOnline(host.NodeCode), 1, 2),
				Domain:             host.Node.Domain,
				CustomDomain:       utils.TrinaryOperation(cache2.GetNodeCustomDomain(host.NodeCode), 1, 2),
				AllowDomainMatcher: host.Node.AllowDomainMatcher,
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
			WhiteEnable: host.WhiteEnable,
			WhiteList:   host.GetWhiteList(),
		})
	}
	return list, total
}
