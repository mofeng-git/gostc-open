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
		host, _ := tx.GostClientHost.Preload(tx.GostClientHost.Node).Where(tx.GostClientHost.Code.Eq(req.Code)).First()
		if host == nil {
			return nil
		}
		if host.CustomDomain != "" {
			_, _ = tx.GostClientHostDomain.Where(tx.GostClientHostDomain.Domain.Eq(host.CustomDomain)).Delete()
		}
		_, _ = tx.GostNodeDomain.Where(
			tx.GostNodeDomain.Prefix.Eq(host.DomainPrefix),
			tx.GostNodeDomain.NodeCode.Eq(host.NodeCode),
		).Delete()
		_, _ = tx.GostAuth.Where(tx.GostAuth.TunnelCode.Eq(host.Code)).Delete()

		if _, err := tx.GostClientHost.Where(tx.GostClientHost.Code.Eq(host.Code)).Delete(); err != nil {
			log.Error("删除用户域名解析失败", zap.Error(err))
			return errors.New("操作失败")
		}
		gost_engine.ClientRemoveHostConfig(tx, *host, host.Node)
		cache.DelTunnelInfo(req.Code)
		return nil
	})
}
