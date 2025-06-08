package service

import (
	"errors"
	"go.uber.org/zap"
	"server/pkg/jwt"
	"server/repository"
	"server/repository/query"
	"server/service/common/cache"
	"server/service/engine"
)

type DeleteReq struct {
	Code string `binding:"required" json:"code" label:"编号"`
}

func (service *service) Delete(claims jwt.Claims, req DeleteReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		proxy, _ := tx.GostClientProxy.Where(
			tx.GostClientProxy.Code.Eq(req.Code),
			tx.GostClientProxy.UserCode.Eq(claims.Code),
		).First()
		if proxy == nil {
			return nil
		}
		_, _ = tx.GostAuth.Where(tx.GostAuth.TunnelCode.Eq(proxy.Code)).Delete()
		if _, err := tx.GostClientProxy.Where(tx.GostClientProxy.Code.Eq(proxy.Code)).Delete(); err != nil {
			log.Error("删除用户代理隧道失败", zap.Error(err))
			return errors.New("操作失败")
		}
		engine.ClientRemoveProxyConfig(*proxy)
		cache.DelTunnelInfo(req.Code)
		return nil
	})
}
