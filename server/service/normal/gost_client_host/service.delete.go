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
		host, _ := tx.GostClientHost.Preload(tx.GostClientHost.Node).Where(
			tx.GostClientHost.UserCode.Eq(claims.Code),
			tx.GostClientHost.Code.Eq(req.Code),
		).First()
		if host == nil {
			return errors.New("操作失败")
		}
		if host.CustomDomain != "" {
			_, _ = tx.GostClientHostDomain.Where(tx.GostClientHostDomain.Domain.Eq(host.CustomDomain)).Delete()
		}
		_, _ = tx.GostNodeDomain.Where(
			tx.GostNodeDomain.Prefix.Eq(host.DomainPrefix),
			tx.GostNodeDomain.NodeCode.Eq(host.NodeCode),
		).Delete()
		if _, err := tx.GostClientHost.Where(tx.GostClientHost.Code.Eq(host.Code)).Delete(); err != nil {
			log.Error("删除用户域名解析失败", zap.Error(err))
			return errors.New("操作失败")
		}
		engine.ClientRemoveHostConfig(tx, *host, host.Node)
		cache.DelTunnelInfo(req.Code)
		cache.DelAdmissionInfo(req.Code)
		return nil
	})
}
