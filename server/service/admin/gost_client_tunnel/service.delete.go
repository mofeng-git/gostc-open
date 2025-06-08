package service

import (
	"errors"
	"go.uber.org/zap"
	"server/repository"
	"server/repository/query"
	"server/service/common/cache"
	"server/service/engine"
)

type DeleteReq struct {
	Code string `binding:"required" json:"code" label:"编号"`
}

func (service *service) Delete(req DeleteReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		tunnel, _ := tx.GostClientTunnel.Preload(tx.GostClientTunnel.Node).Where(tx.GostClientTunnel.Code.Eq(req.Code)).First()
		if tunnel == nil {
			return nil
		}
		_, _ = tx.GostAuth.Where(tx.GostAuth.TunnelCode.Eq(tunnel.Code)).Delete()
		if _, err := tx.GostClientTunnel.Where(tx.GostClientTunnel.Code.Eq(tunnel.Code)).Delete(); err != nil {
			log.Error("删除用户私有隧道失败", zap.Error(err))
			return errors.New("操作失败")
		}
		engine.ClientRemoveTunnelConfig(tx, *tunnel, tunnel.Node)
		cache.DelTunnelInfo(req.Code)
		return nil
	})
}
