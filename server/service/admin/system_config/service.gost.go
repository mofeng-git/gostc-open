package service

import (
	"errors"
	"go.uber.org/zap"
	"server/model"
	"server/repository"
	"server/repository/query"
	"server/service/common/cache"
)

type GostReq struct {
	Version     string `binding:"required" json:"version"`
	FuncWeb     string `binding:"required" json:"funcWeb"`
	FuncForward string `binding:"required" json:"funcForward"`
	FuncTunnel  string `binding:"required" json:"funcTunnel"`
	FuncP2P     string `binding:"required" json:"funcP2P"`
	FuncProxy   string `binding:"required" json:"funcProxy"`
	FuncTun     string `binding:"required" json:"funcTun"`
	FuncNode    string `binding:"required" json:"funcNode"`
}

func (service *service) Gost(req GostReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		_, _ = tx.SystemConfig.Where(tx.SystemConfig.Kind.Eq(model.SYSTEM_CONFIG_KIND_GOST)).Delete()
		if err := tx.SystemConfig.Create(model.GenerateSystemConfigGost(
			req.Version,
			"2",
			req.FuncWeb,
			req.FuncForward,
			req.FuncTunnel,
			req.FuncP2P,
			req.FuncProxy,
			req.FuncTun,
			req.FuncNode,
		)...); err != nil {
			log.Error("修改系统GOST配置失败", zap.Error(err))
			return errors.New("操作失败")
		}
		cache.SetSystemConfigGost(model.SystemConfigGost{
			Version:     req.Version,
			Logger:      "2",
			FuncWeb:     req.FuncWeb,
			FuncForward: req.FuncForward,
			FuncTunnel:  req.FuncTunnel,
			FuncP2P:     req.FuncP2P,
			FuncProxy:   req.FuncProxy,
			FuncTun:     req.FuncTun,
			FuncNode:    req.FuncNode,
		})
		return nil
	})
}
