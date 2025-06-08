package service

import (
	"errors"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"server/model"
	"server/repository"
	"server/repository/query"
	"server/service/common/cache"
	"server/service/engine"
	"time"
)

type ConfigReq struct {
	Code         string `binding:"required" json:"code"`
	ChargingType int    `json:"chargingType"`
	Cycle        int    `json:"cycle"`
	Amount       string `json:"amount"`
	Limiter      int    `json:"limiter"`
	RLimiter     int    `json:"rLimiter"`
	CLimiter     int    `json:"cLimiter"`
	ExpAt        string `json:"expAt"`
}

func (service *service) Config(req ConfigReq) error {
	db, _, log := repository.Get("")
	expAt, err := time.ParseInLocation(time.DateTime, req.ExpAt, time.Local)
	if err != nil {
		return errors.New("到期时间错误")
	}

	return db.Transaction(func(tx *query.Query) error {
		proxy, _ := tx.GostClientProxy.Where(tx.GostClientProxy.Code.Eq(req.Code)).First()
		if proxy == nil {
			return nil
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
		case model.GOST_CONFIG_CHARGING_FREE:
		default:
			return errors.New("计费类型错误")
		}

		proxy.ChargingType = req.ChargingType
		proxy.Cycle = req.Cycle
		proxy.Amount = amount
		proxy.Limiter = req.Limiter
		//proxy.RLimiter = req.RLimiter
		//proxy.CLimiter = req.CLimiter
		proxy.ExpAt = expAt.Unix()
		if err = tx.GostClientProxy.Save(proxy); err != nil {
			log.Error("修改代理隧道配置失败", zap.Error(err))
			return errors.New("操作失败")
		}
		engine.ClientProxyConfig(tx, proxy.Code)
		cache.SetTunnelInfo(cache.TunnelInfo{
			Code:        proxy.Code,
			Type:        model.GOST_TUNNEL_TYPE_PROXY,
			ClientCode:  proxy.ClientCode,
			UserCode:    proxy.UserCode,
			NodeCode:    proxy.NodeCode,
			ChargingTye: proxy.ChargingType,
			ExpAt:       proxy.ExpAt,
			Limiter:     proxy.Limiter,
		})
		return nil
	})
}
