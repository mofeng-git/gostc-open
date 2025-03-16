package service

import (
	"errors"
	"go.uber.org/zap"
	"server/repository"
	"server/repository/query"
	"server/service/common/cache"
	"server/service/gost_engine"
)

type DeleteReq struct {
	Code string `binding:"required" json:"code" label:"编号"`
}

func (service *service) Delete(req DeleteReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		proxy, _ := tx.GostClientProxy.Preload(tx.GostClientProxy.Node).Where(tx.GostClientProxy.Code.Eq(req.Code)).First()
		if proxy == nil {
			return nil
		}
		_, _ = tx.GostAuth.Where(tx.GostAuth.TunnelCode.Eq(proxy.Code)).Delete()
		if _, err := tx.GostClientProxy.Where(tx.GostClientProxy.Code.Eq(proxy.Code)).Delete(); err != nil {
			log.Error("删除用户代理隧道失败", zap.Error(err))
			return errors.New("操作失败")
		}
		gost_engine.ClientRemoveProxyConfig(*proxy, proxy.Node)
		cache.DelTunnelInfo(req.Code)
		return nil
	})
}
