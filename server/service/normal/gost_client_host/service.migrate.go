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
		host, _ := tx.GostClientHost.Preload(
			tx.GostClientHost.Node,
		).Where(
			tx.GostClientHost.UserCode.Eq(user.Code),
			tx.GostClientHost.Code.Eq(req.Code),
		).First()
		if host == nil {
			return errors.New("操作失败")
		}
		if host.ClientCode == req.ClientCode {
			return nil
		}
		if host.Enable == 1 {
			gost_engine.ClientRemoveHostConfig(tx, *host, host.Node)
		}
		host.ClientCode = req.ClientCode
		if err := tx.GostClientHost.Save(host); err != nil {
			log.Error("迁移域名解析失败", zap.Error(err))
			return errors.New("操作失败")
		}
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
		gost_engine.ClientHostConfig(tx, host.Code)
		return nil
	})
}
