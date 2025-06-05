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
		p2p, _ := tx.GostClientP2P.Where(tx.GostClientP2P.Code.Eq(req.Code)).First()
		if p2p == nil {
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

		p2p.ChargingType = req.ChargingType
		p2p.Cycle = req.Cycle
		p2p.Amount = amount
		p2p.Limiter = req.Limiter
		p2p.RLimiter = req.RLimiter
		p2p.CLimiter = req.CLimiter
		p2p.ExpAt = expAt.Unix()
		if err = tx.GostClientP2P.Save(p2p); err != nil {
			log.Error("修改P2P隧道配置失败", zap.Error(err))
			return errors.New("操作失败")
		}
		gost_engine.ClientP2PConfig(tx, p2p.Code)
		cache.SetTunnelInfo(cache.TunnelInfo{
			Code:        p2p.Code,
			Type:        model.GOST_TUNNEL_TYPE_P2P,
			ClientCode:  p2p.ClientCode,
			UserCode:    p2p.UserCode,
			NodeCode:    p2p.NodeCode,
			ChargingTye: p2p.ChargingType,
			ExpAt:       p2p.ExpAt,
			Limiter:     p2p.Limiter,
		})
		return nil
	})
}
