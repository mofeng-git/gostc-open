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

		proxy, _ := tx.GostClientProxy.Where(
			tx.GostClientProxy.UserCode.Eq(user.Code),
			tx.GostClientProxy.Code.Eq(req.Code),
		).First()
		if proxy == nil {
			return errors.New("操作失败")
		}

		var expAt = time.Unix(proxy.ExpAt, 0)
		if expAt.Unix() < time.Now().Unix() {
			expAt = time.Now()
		}
		switch proxy.ChargingType {
		case model.GOST_CONFIG_CHARGING_CUCLE_DAY:
			expAt = expAt.Add(time.Duration(proxy.Cycle) * 24 * time.Hour)
			if user.Amount.LessThan(proxy.Amount) {
				return errors.New("积分不足")
			}

			if _, err := tx.SystemUser.Where(
				tx.SystemUser.Code.Eq(user.Code),
				tx.SystemUser.Version.Eq(user.Version),
			).UpdateSimple(
				tx.SystemUser.Amount.Value(user.Amount.Sub(proxy.Amount)),
				tx.SystemUser.Version.Value(user.Version+1),
			); err != nil {
				log.Error("扣减积分失败", zap.Error(err))
				return errors.New("操作失败")
			}
		case model.GOST_CONFIG_CHARGING_ONLY_ONCE, model.GOST_CONFIG_CHARGING_FREE:
			return nil
		}
		proxy.Status = 1
		proxy.ExpAt = expAt.Unix()
		if err := tx.GostClientProxy.Save(proxy); err != nil {
			log.Error("续费用户端口转发失败", zap.Error(err))
			return errors.New("操作失败")
		}
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
		engine.ClientProxyConfig(tx, proxy.Code)
		return nil
	})
}
