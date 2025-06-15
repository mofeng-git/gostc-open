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

		host, _ := tx.GostClientHost.Where(
			tx.GostClientHost.UserCode.Eq(user.Code),
			tx.GostClientHost.Code.Eq(req.Code),
		).First()
		if host == nil {
			return errors.New("操作失败")
		}

		var expAt = time.Unix(host.ExpAt, 0)
		if expAt.Unix() < time.Now().Unix() {
			expAt = time.Now()
		}
		switch host.ChargingType {
		case model.GOST_CONFIG_CHARGING_CUCLE_DAY:
			expAt = expAt.Add(time.Duration(host.Cycle) * 24 * time.Hour)
			if user.Amount.LessThan(host.Amount) {
				return errors.New("积分不足")
			}

			if _, err := tx.SystemUser.Where(
				tx.SystemUser.Code.Eq(user.Code),
				tx.SystemUser.Version.Eq(user.Version),
			).UpdateSimple(
				tx.SystemUser.Amount.Value(user.Amount.Sub(host.Amount)),
				tx.SystemUser.Version.Value(user.Version+1),
			); err != nil {
				log.Error("扣减积分失败", zap.Error(err))
				return errors.New("操作失败")
			}
		case model.GOST_CONFIG_CHARGING_ONLY_ONCE, model.GOST_CONFIG_CHARGING_FREE:
			return nil
		}
		host.Status = 1
		host.ExpAt = expAt.Unix()
		if err := tx.GostClientHost.Save(host); err != nil {
			log.Error("续费用户端口转发失败", zap.Error(err))
			return errors.New("操作失败")
		}
		engine.ClientHostConfig(tx, host.Code)
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
