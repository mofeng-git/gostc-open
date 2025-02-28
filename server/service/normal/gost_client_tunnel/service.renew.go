package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/model"
	"server/pkg/jwt"
	"server/repository"
	"server/service/common/cache"
	"server/service/gost_engine"
	"time"
)

type RenewReq struct {
	Code string `binding:"required" json:"code"`
}

func (service *service) Renew(claims jwt.Claims, req RenewReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *gorm.DB) error {
		var user model.SystemUser
		if tx.Where("code = ?", claims.Code).First(&user).RowsAffected == 0 {
			return errors.New("用户错误")
		}

		var tunnel model.GostClientTunnel
		if tx.Where("code = ? AND user_code = ?", req.Code, user.Code).First(&tunnel).RowsAffected == 0 {
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
			user.Amount = user.Amount.Sub(tunnel.Amount)
			if err := tx.Save(&user).Error; err != nil {
				log.Error("扣减积分失败", zap.Error(err))
				return errors.New("操作失败")
			}
		case model.GOST_CONFIG_CHARGING_ONLY_ONCE, model.GOST_CONFIG_CHARGING_FREE:
			return nil
		}
		tunnel.Status = 1
		tunnel.ExpAt = expAt.Unix()
		if err := tx.Save(&tunnel).Error; err != nil {
			log.Error("续费用户端口转发失败", zap.Error(err))
			return errors.New("操作失败")
		}
		gost_engine.ClientForwardConfig(tx, tunnel.Code)
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
