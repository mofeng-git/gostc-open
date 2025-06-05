package service

import (
	"errors"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"server/model"
	"server/repository"
	"server/repository/query"
	"server/service/common/cache"
	"server/service/gost_engine"
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
		tunnel, _ := tx.GostClientTunnel.Where(tx.GostClientTunnel.Code.Eq(req.Code)).First()
		if tunnel == nil {
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

		tunnel.ChargingType = req.ChargingType
		tunnel.Cycle = req.Cycle
		tunnel.Amount = amount
		tunnel.Limiter = req.Limiter
		tunnel.RLimiter = req.RLimiter
		tunnel.CLimiter = req.CLimiter
		tunnel.ExpAt = expAt.Unix()
		if err = tx.GostClientTunnel.Save(tunnel); err != nil {
			log.Error("修改私有隧道配置失败", zap.Error(err))
			return errors.New("操作失败")
		}
		gost_engine.ClientTunnelConfig(tx, tunnel.Code)
		cache.SetTunnelInfo(cache.TunnelInfo{
			Code:        tunnel.Code,
			Type:        model.GOST_TUNNEL_TYPE_TUNNEL,
			ClientCode:  tunnel.ClientCode,
			UserCode:    tunnel.UserCode,
			NodeCode:    tunnel.NodeCode,
			ChargingTye: tunnel.ChargingType,
			ExpAt:       tunnel.ExpAt,
			Limiter:     tunnel.Limiter,
		})
		return nil
	})
}
