package service

import (
	"server/model"
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
	UserAccount   string        `json:"userAccount"`
	Code          string        `json:"code"`
	Name          string        `json:"name"`
	TargetIp      string        `json:"targetIp"`
	TargetPort    string        `json:"targetPort"`
	Port          string        `json:"port"`
	Node          ItemNode      `json:"node"`
	Client        ItemClient    `json:"client"`
	Config        ItemConfig    `json:"config"`
	Enable        int           `json:"enable"`
	WarnMsg       string        `json:"warnMsg"`
	CreatedAt     string        `json:"createdAt"`
	InputBytes    int64         `json:"inputBytes"`
	OutputBytes   int64         `json:"outputBytes"`
	MatcherEnable int           `json:"matcherEnable"`
	Matchers      []ItemMatcher `json:"matchers"`
	TcpMatcher    ItemMatcher   `json:"tcpMatcher"`
	SSHMatcher    ItemMatcher   `json:"sshMatcher"`
}

type ItemMatcher struct {
	Host       string `json:"host"`
	TargetIp   string `json:"targetIp"`
	TargetPort string `json:"targetPort"`
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

func (service *service) Page(req PageReq) (list []Item, total int64) {
	db, _, _ := repository.Get("")
	var forwards []model.GostClientForward
	var where = db
	if req.Account != "" {
		where = where.Where(
			"user_code in (?)",
			db.Model(&model.SystemUser{}).Where("account like ?", "%"+req.Account+"%").Select("code"),
		)
	}
	if req.Enable > 0 {
		where = where.Where("enable = ?", req.Enable)
	}
	db.Model(&forwards).Count(&total)
	db.
		Preload("User").
		Preload("Client").
		Preload("Node").
		Where(where).Order("id desc").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&forwards)
	for _, forward := range forwards {
		var mathcers []ItemMatcher
		for _, item := range forward.GetMatcher() {
			mathcers = append(mathcers, ItemMatcher{
				Host:       item.Host,
				TargetIp:   item.TargetIp,
				TargetPort: item.TargetPort,
			})
		}
		tcpIp, tcpPort := forward.GetTcpMatcher()
		sshIp, sshPort := forward.GetSSHMatcher()
		obsInfo := cache.GetTunnelObsDateRange(cache.MONTH_DATEONLY_LIST, forward.Code)
		list = append(list, Item{
			UserAccount: forward.User.Account,
			Code:        forward.Code,
			Name:        forward.Name,
			TargetIp:    forward.TargetIp,
			TargetPort:  forward.TargetPort,
			Port:        forward.Port,
			Node: ItemNode{
				Code:    forward.NodeCode,
				Name:    forward.Node.Name,
				Address: forward.Node.Address,
				Online:  utils.TrinaryOperation(cache.GetNodeOnline(forward.NodeCode), 1, 2),
			},
			Client: ItemClient{
				Code:   forward.ClientCode,
				Name:   forward.Client.Name,
				Online: utils.TrinaryOperation(cache.GetClientOnline(forward.ClientCode), 1, 2),
			},
			Config: ItemConfig{
				ChargingType: forward.ChargingType,
				Cycle:        forward.Cycle,
				Amount:       forward.Amount.String(),
				Limiter:      forward.Limiter,
				RLimiter:     forward.RLimiter,
				CLimiter:     forward.CLimiter,
				ExpAt:        time.Unix(forward.ExpAt, 0).Format(time.DateTime),
				OnlyChina:    forward.OnlyChina,
			},
			Enable:        forward.Enable,
			WarnMsg:       warn_msg.GetForwardWarnMsg(forward),
			CreatedAt:     forward.CreatedAt.Format(time.DateTime),
			InputBytes:    obsInfo.InputBytes,
			OutputBytes:   obsInfo.OutputBytes,
			MatcherEnable: forward.MatcherEnable,
			Matchers:      mathcers,
			TcpMatcher: ItemMatcher{
				TargetIp:   tcpIp,
				TargetPort: tcpPort,
			},
			SSHMatcher: ItemMatcher{
				TargetIp:   sshIp,
				TargetPort: sshPort,
			},
		})
	}
	return list, total
}
