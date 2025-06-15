package service

import (
	"errors"
	"go.uber.org/zap"
	"server/pkg/jwt"
	"server/repository"
	"server/repository/cache"
	"server/repository/query"
	"server/service/engine"
)

type DeleteReq struct {
	Code string `binding:"required" json:"code" label:"编号"`
}

func (service *service) Delete(claims jwt.Claims, req DeleteReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		p2p, _ := tx.GostClientP2P.Where(tx.GostClientP2P.Code.Eq(req.Code), tx.GostClientP2P.UserCode.Eq(claims.Code)).First()
		if p2p == nil {
			return nil
		}
		_, _ = tx.GostAuth.Where(tx.GostAuth.TunnelCode.Eq(p2p.Code)).Delete()
		if _, err := tx.GostClientP2P.Where(tx.GostClientP2P.Code.Eq(p2p.Code)).Delete(); err != nil {
			log.Error("删除用户P2P隧道失败", zap.Error(err))
			return errors.New("操作失败")
		}
		engine.ClientRemoveP2PConfig(*p2p)
		cache.DelTunnelInfo(req.Code)
		return nil
	})
}
