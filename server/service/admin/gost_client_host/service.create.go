package service

import (
	"errors"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"server/model"
	"server/pkg/utils"
	"server/repository"
	"server/repository/query"
	"server/service/common/cache"
	"server/service/gost_engine"
	"time"
)

type CreateReq struct {
	UserCode   string `binding:"required" json:"userCode" label:"用户编号"`
	Name       string `binding:"required" json:"name" label:"名称"`
	TargetIp   string `binding:"required" json:"targetIp" label:"内网IP"`
	TargetPort string `binding:"required" json:"targetPort" label:"内网端口"`
	NodeCode   string `binding:"required" json:"nodeCode" label:"节点编号"`
	ClientCode string `binding:"required" json:"clientCode" label:"客户端编号"`

	ChargingType int    `json:"chargingType"`
	Cycle        int    `json:"cycle"`
	Amount       string `json:"amount"`
	Limiter      int    `json:"limiter"`
	RLimiter     int    `json:"rLimiter"`
	CLimiter     int    `json:"cLimiter"`
	OnlyChina    int    `json:"onlyChina"`
	ExpAt        string `json:"expAt"`
}

func (service *service) Create(req CreateReq) error {
	db, _, log := repository.Get("")
	if !utils.ValidateLocalIP(req.TargetIp) {
		return errors.New("内网IP格式错误")
	}
	if !utils.ValidatePort(req.TargetPort) {
		return errors.New("内网端口格式错误")
	}
	expAt, err := time.ParseInLocation(time.DateTime, req.ExpAt, time.Local)
	if err != nil {
		return errors.New("到期时间错误")
	}

	return db.Transaction(func(tx *query.Query) error {
		user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(req.UserCode)).First()
		if user == nil {
			return errors.New("用户错误")
		}

		var domainPrefix = utils.RandStr(8, utils.LatterDict)
		node, _ := tx.GostNode.Where(tx.GostNode.Code.Eq(req.NodeCode)).First()
		if node == nil {
			return errors.New("节点错误")
		}
		if !node.CheckDomainPrefix(domainPrefix) {
			return errors.New("该域名前缀被禁止使用")
		}
		if node.Web != 1 {
			return errors.New("该节点未启用域名解析功能")
		}

		if err := tx.GostNodeDomain.Create(&model.GostNodeDomain{
			Prefix:   domainPrefix,
			NodeCode: node.Code,
		}); err != nil {
			return errors.New("该域名前缀已被使用")
		}

		client, _ := tx.GostClient.Where(
			tx.GostClient.Code.Eq(req.ClientCode),
			tx.GostClient.UserCode.Eq(req.UserCode),
		).First()
		if client == nil {
			return errors.New("客户端错误")
		}

		var amount decimal.Decimal
		switch req.ChargingType {
		case model.GOST_CONFIG_CHARGING_CUCLE_DAY, model.GOST_CONFIG_CHARGING_ONLY_ONCE:
			amount, err = decimal.NewFromString(req.Amount)
			if err != nil {
				return errors.New("套餐积分错误")
			}
			if req.ChargingType == model.GOST_CONFIG_CHARGING_CUCLE_DAY && req.Cycle <= 0 {
				return errors.New("计费循环周期错误")
			}
		default:
			return errors.New("计费类型错误")
		}

		var host = model.GostClientHost{
			Name:         req.Name,
			TargetIp:     req.TargetIp,
			TargetPort:   req.TargetPort,
			DomainPrefix: domainPrefix,
			NodeCode:     req.NodeCode,
			ClientCode:   req.ClientCode,
			UserCode:     req.UserCode,
			GostClientConfig: model.GostClientConfig{
				ChargingType: req.ChargingType,
				Cycle:        req.Cycle,
				Amount:       amount,
				Limiter:      req.Limiter,
				RLimiter:     req.RLimiter,
				CLimiter:     req.CLimiter,
				OnlyChina:    req.OnlyChina,
				ExpAt:        expAt.Unix(),
			},
		}
		if err := tx.GostClientHost.Save(&host); err != nil {
			log.Error("新增用户域名解析失败", zap.Error(err))
			return errors.New("操作失败")
		}
		var auth = model.GostAuth{
			TunnelType: model.GOST_TUNNEL_TYPE_FORWARD,
			TunnelCode: host.Code,
			User:       utils.RandStr(10, utils.AllDict),
			Password:   utils.RandStr(10, utils.AllDict),
		}
		if err = tx.GostAuth.Create(&auth); err != nil {
			log.Error("生成授权信息失败", zap.Error(err))
			return errors.New("操作失败")
		}
		cache.SetGostAuth(auth.User, auth.Password, host.Code)
		gost_engine.ClientHostConfig(tx, host.Code)
		cache.SetTunnelInfo(cache.TunnelInfo{
			Code:        host.Code,
			Type:        model.GOST_TUNNEL_TYPE_HOST,
			ClientCode:  host.ClientCode,
			UserCode:    host.UserCode,
			NodeCode:    host.NodeCode,
			ChargingTye: host.ChargingType,
			ExpAt:       host.ExpAt,
			Limiter:     host.Limiter,
		})
		return nil
	})
}
