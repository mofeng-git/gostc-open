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
		tunnel, _ := tx.GostClientTunnel.Where(tx.GostClientTunnel.Code.Eq(req.Code), tx.GostClientTunnel.UserCode.Eq(claims.Code)).First()
		if tunnel == nil {
			return errors.New("操作失败")
		}

		var expAt = time.Unix(tunnel.ExpAt, 0)
		if expAt.Unix() < time.Now().Unix() {
			expAt = time.Now()
		}
		switch tunnel.ChargingType {
		case model.GOST_CONFIG_CHARGING_CUCLE_DAY:
			expAt = expAt.Add(time.Duration(tunnel.Cycle) * 24 * time.Hour)
			if user.Amount.LessThan(tunnel.Amount) {
				return errors.New("积分不足")
			}

			if _, err := tx.SystemUser.Where(
				tx.SystemUser.Code.Eq(user.Code),
				tx.SystemUser.Version.Eq(user.Version),
			).UpdateSimple(
				tx.SystemUser.Amount.Value(user.Amount.Sub(tunnel.Amount)),
				tx.SystemUser.Version.Value(user.Version+1),
			); err != nil {
				log.Error("扣减积分失败", zap.Error(err))
				return errors.New("操作失败")
			}
		case model.GOST_CONFIG_CHARGING_ONLY_ONCE, model.GOST_CONFIG_CHARGING_FREE:
			return nil
		}
		tunnel.Status = 1
		tunnel.ExpAt = expAt.Unix()
		if err := tx.GostClientTunnel.Save(tunnel); err != nil {
			log.Error("续费用户端口转发失败", zap.Error(err))
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
