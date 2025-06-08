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
	UserAccount   string        `json:"userAccount"`
	Code          string        `json:"code"`
	Name          string        `json:"name"`
	TargetIp      string        `json:"targetIp"`
	TargetPort    string        `json:"targetPort"`
	ProxyProtocol int           `json:"proxyProtocol"`
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
	ExpAt        string `json:"expAt"`
}

func (service *service) Page(req PageReq) (list []Item, total int64) {
	db, _, _ := repository.Get("")
	var where []gen.Condition
	if req.Account != "" {
		var userCodes []string
		_ = db.SystemUser.Where(db.SystemUser.Account.Like("%"+req.Account+"%")).Pluck(db.SystemUser.Code, &userCodes)
		where = append(where, db.GostClientForward.UserCode.In(userCodes...))
	}
	if req.NodeName != "" {
		var codes []string
		_ = db.GostNode.Where(db.GostNode.Name.Like("%"+req.NodeName+"%")).Pluck(db.GostNode.Code, &codes)
		where = append(where, db.GostClientForward.NodeCode.In(codes...))
	}
	if req.Enable > 0 {
		where = append(where, db.GostClientForward.Enable.Eq(req.Enable))
	}
	if req.Name != "" {
		where = append(where, db.GostClientForward.Name.Like("%"+req.Name+"%"))
	}
	if req.ClientName != "" {
		var clientCodes []string
		_ = db.GostClient.Where(db.GostClient.Name.Like("%"+req.ClientName+"%")).Pluck(db.GostClient.Code, &clientCodes)
		where = append(where, db.GostClientForward.ClientCode.In(clientCodes...))
	}

	forwards, total, _ := db.GostClientForward.Preload(
		db.GostClientForward.User,
		db.GostClientForward.Client,
		db.GostClientForward.Node,
	).Where(where...).Order(db.GostClientForward.Id.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	for _, forward := range forwards {
		//var mathcers []ItemMatcher
		//for _, item := range forward.GetMatcher() {
		//	mathcers = append(mathcers, ItemMatcher{
		//		Host:       item.Host,
		//		TargetIp:   item.TargetIp,
		//		TargetPort: item.TargetPort,
		//	})
		//}
		//tcpIp, tcpPort := forward.GetTcpMatcher()
		//sshIp, sshPort := forward.GetSSHMatcher()
		obsInfo := cache.GetTunnelObsDateRange(cache.MONTH_DATEONLY_LIST, forward.Code)
		list = append(list, Item{
			UserAccount:   forward.User.Account,
			Code:          forward.Code,
			Name:          forward.Name,
			TargetIp:      forward.TargetIp,
			TargetPort:    forward.TargetPort,
			ProxyProtocol: forward.ProxyProtocol,
			Port:          forward.Port,
			Node: ItemNode{
				Code: forward.NodeCode,
				Name: forward.Node.Name,
				Address: func() string {
					address, _, _ := net.SplitHostPort(forward.Node.Address)
					return address
				}(),
				Online: utils.TrinaryOperation(cache.GetNodeOnline(forward.NodeCode), 1, 2),
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
				//RLimiter:     forward.RLimiter,
				//CLimiter:     forward.CLimiter,
				ExpAt: time.Unix(forward.ExpAt, 0).Format(time.DateTime),
			},
			Enable:      forward.Enable,
			WarnMsg:     warn_msg.GetForwardWarnMsg(*forward),
			CreatedAt:   forward.CreatedAt.Format(time.DateTime),
			InputBytes:  obsInfo.InputBytes,
			OutputBytes: obsInfo.OutputBytes,
			//MatcherEnable: forward.MatcherEnable,
			//Matchers:      mathcers,
			//TcpMatcher: ItemMatcher{
			//	TargetIp:   tcpIp,
			//	TargetPort: tcpPort,
			//},
			//SSHMatcher: ItemMatcher{
			//	TargetIp:   sshIp,
			//	TargetPort: sshPort,
			//},
		})
	}
	return list, total
}
