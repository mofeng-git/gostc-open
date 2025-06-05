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
		p2p, _ := tx.GostClientP2P.Where(tx.GostClientP2P.Code.Eq(req.Code), tx.GostClientP2P.UserCode.Eq(claims.Code)).First()
		if p2p == nil {
			return errors.New("操作失败")
		}
		if p2p.ClientCode == req.ClientCode {
			return nil
		}
		if p2p.Enable == 1 {
			gost_engine.ClientRemoveP2PConfig(*p2p)
		}
		p2p.ClientCode = req.ClientCode
		if err := tx.GostClientP2P.Save(p2p); err != nil {
			log.Error("迁移P2P隧道失败", zap.Error(err))
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
		gost_engine.ClientP2PConfig(tx, p2p.Code)
		return nil
	})
}
