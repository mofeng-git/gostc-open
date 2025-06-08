package service

import (
	"errors"
	"go.uber.org/zap"
	"server/model"
	"server/pkg/jwt"
	"server/repository"
	"server/repository/query"
	"server/service/common/cache"
	"server/service/engine"
)

type MigrateReq struct {
	Code       string `binding:"required" json:"code"`
	ClientCode string `binding:"required" json:"clientCode"`
}

func (service *service) Migrate(claims jwt.Claims, req MigrateReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(claims.Code)).First()
		if user == nil {
			return errors.New("用户错误")
		}
		newClient, _ := tx.GostClient.Where(tx.GostClient.UserCode.Eq(claims.Code), tx.GostClient.Code.Eq(req.ClientCode)).First()
		if newClient == nil {
			return errors.New("新客户端不存在")
		}

		forward, _ := tx.GostClientForward.Where(
			tx.GostClientForward.UserCode.Eq(user.Code),
			tx.GostClientForward.Code.Eq(req.Code),
		).First()
		if forward == nil {
			return errors.New("操作失败")
		}
		if forward.ClientCode == req.ClientCode {
			return nil
		}
		if forward.Enable == 1 {
			engine.ClientRemoveForwardConfig(*forward)
		}

		forward.ClientCode = req.ClientCode
		if err := tx.GostClientForward.Save(forward); err != nil {
			log.Error("迁移端口转发失败", zap.Error(err))
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
