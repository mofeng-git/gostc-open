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
		tunnel, _ := tx.GostClientTunnel.Preload(
			tx.GostClientTunnel.Node,
		).Where(
			tx.GostClientTunnel.Code.Eq(req.Code),
			tx.GostClientTunnel.UserCode.Eq(claims.Code),
		).First()
		if tunnel == nil {
			return errors.New("操作失败")
		}

		if tunnel.ClientCode == req.ClientCode {
			return nil
		}

		if tunnel.Enable == 1 {
			gost_engine.ClientRemoveTunnelConfig(tx, *tunnel, tunnel.Node)
		}

		tunnel.ClientCode = req.ClientCode
		if err := tx.GostClientTunnel.Save(tunnel); err != nil {
			log.Error("迁移私有隧道失败", zap.Error(err))
			return errors.New("操作失败")
		}
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
		gost_engine.ClientTunnelConfig(tx, tunnel.Code)
		return nil
	})
}
