package service

import (
	"errors"
	"go.uber.org/zap"
	"server/model"
	"server/pkg/jwt"
	"server/repository"
	"server/repository/cache"
	"server/repository/query"
	"server/service/engine"
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
		p2p, _ := tx.GostClientP2P.Where(tx.GostClientP2P.Code.Eq(req.Code), tx.GostClientP2P.UserCode.Eq(claims.Code)).First()
		if p2p == nil {
			return errors.New("操作失败")
		}

		var expAt = time.Unix(p2p.ExpAt, 0)
		if expAt.Unix() < time.Now().Unix() {
			expAt = time.Now()
		}
		switch p2p.ChargingType {
		case model.GOST_CONFIG_CHARGING_CUCLE_DAY:
			expAt = expAt.Add(time.Duration(p2p.Cycle) * 24 * time.Hour)
			if user.Amount.LessThan(p2p.Amount) {
				return errors.New("积分不足")
			}

			if _, err := tx.SystemUser.Where(
				tx.SystemUser.Code.Eq(user.Code),
				tx.SystemUser.Version.Eq(user.Version),
			).UpdateSimple(
				tx.SystemUser.Amount.Value(user.Amount.Sub(p2p.Amount)),
				tx.SystemUser.Version.Value(user.Version+1),
			); err != nil {
				log.Error("扣减积分失败", zap.Error(err))
				return errors.New("操作失败")
			}
		case model.GOST_CONFIG_CHARGING_ONLY_ONCE, model.GOST_CONFIG_CHARGING_FREE:
			return nil
		}
		p2p.Status = 1
		p2p.ExpAt = expAt.Unix()
		if err := tx.GostClientP2P.Save(p2p); err != nil {
			log.Error("续费用户端口转发失败", zap.Error(err))
			return errors.New("操作失败")
		}
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
		engine.ClientP2PConfig(tx, p2p.Code)
		return nil
	})
}
