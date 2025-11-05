package service

import (
	"errors"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"server/model"
	"server/repository"
	"server/repository/cache"
	"server/repository/query"
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
		forward, _ := tx.GostClientForward.Where(tx.GostClientForward.Code.Eq(req.Code)).First()
		if forward == nil {
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

		forward.ChargingType = req.ChargingType
		forward.Cycle = req.Cycle
		forward.Amount = amount
		forward.Limiter = req.Limiter
		//forward.RLimiter = req.RLimiter
		//forward.CLimiter = req.CLimiter
		forward.ExpAt = expAt.Unix()
		if err = tx.GostClientForward.Save(forward); err != nil {
			log.Error("修改端口转发配置失败", zap.Error(err))
			return errors.New("操作失败")
		}
		cache.SetTunnelInfo(cache.TunnelInfo{
			Code:        forward.Code,
			Type:        model.GOST_TUNNEL_TYPE_FORWARD,
			ClientCode:  forward.ClientCode,
			UserCode:    forward.UserCode,
			NodeCode:    forward.NodeCode,
			ChargingTye: forward.ChargingType,
			ExpAt:       forward.ExpAt,
			Limiter:     forward.Limiter,
		})
		engine.ClientForwardConfig(tx, forward.Code)
		return nil
	})
}
