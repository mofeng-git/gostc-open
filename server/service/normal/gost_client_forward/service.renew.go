package service

import (
	"errors"
	"go.uber.org/zap"
	"server/model"
	"server/pkg/jwt"
	"server/repository"
	"server/repository/query"
	"server/service/common/cache"
	"server/service/gost_engine"
	"time"
)

type RenewReq struct {
	Code string `binding:"required" json:"code"`
}

func (service *service) Renew(claims jwt.Claims, req RenewReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(claims.Code)).First()
		if user == nil {
			return errors.New("用户错误")
		}

		forward, _ := tx.GostClientForward.Where(
			tx.GostClientForward.UserCode.Eq(user.Code),
			tx.GostClientForward.Code.Eq(req.Code),
		).First()
		if forward == nil {
			return errors.New("操作失败")
		}

		var expAt = time.Unix(forward.ExpAt, 0)
		if expAt.Unix() < time.Now().Unix() {
			expAt = time.Now()
		}
		switch forward.ChargingType {
		case model.GOST_CONFIG_CHARGING_CUCLE_DAY:
			expAt = expAt.Add(time.Duration(forward.Cycle) * 24 * time.Hour)
			if user.Amount.LessThan(forward.Amount) {
				return errors.New("积分不足")
			}
			user.Amount = user.Amount.Sub(forward.Amount)
			if err := tx.SystemUser.Save(user); err != nil {
				log.Error("扣减积分失败", zap.Error(err))
				return errors.New("操作失败")
			}
		case model.GOST_CONFIG_CHARGING_ONLY_ONCE, model.GOST_CONFIG_CHARGING_FREE:
			return nil
		}
		forward.Status = 1
		forward.ExpAt = expAt.Unix()
		if err := tx.GostClientForward.Save(forward); err != nil {
			log.Error("续费用户端口转发失败", zap.Error(err))
			return errors.New("操作失败")
		}
		gost_engine.ClientForwardConfig(tx, forward.Code)
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
		return nil
	})
}
