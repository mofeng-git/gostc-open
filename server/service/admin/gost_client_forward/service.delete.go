package service

import (
	"errors"
	"go.uber.org/zap"
	"server/repository"
	"server/repository/cache"
	"server/repository/query"
	"server/service/common/node_port"
	"server/service/engine"
)

type DeleteReq struct {
	Code string `binding:"required" json:"code" label:"编号"`
}

func (service *service) Delete(req DeleteReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		forward, _ := tx.GostClientForward.Preload(tx.GostClientForward.Node).Where(tx.GostClientForward.Code.Eq(req.Code)).First()
		if forward == nil {
			return nil
		}

		_, _ = tx.GostNodePort.Where(tx.GostNodePort.Port.Eq(forward.Port), tx.GostNodePort.NodeCode.Eq(forward.NodeCode)).Delete()
		_, _ = tx.GostAuth.Where(tx.GostAuth.TunnelCode.Eq(forward.Code)).Delete()
		if _, err := tx.GostClientForward.Where(tx.GostClientForward.Code.Eq(forward.Code)).Delete(); err != nil {
			log.Error("删除用户端口转发失败", zap.Error(err))
			return errors.New("操作失败")
		}
		node_port.ReleasePort(forward.NodeCode, forward.Port)
		engine.ClientRemoveForwardConfig(*forward)
		cache.DelTunnelInfo(req.Code)
		cache.DelAdmissionInfo(req.Code)
		return nil
	})
}
