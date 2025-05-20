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
		node, _ := tx.GostNode.Where(tx.GostNode.Code.Eq(req.Code)).First()
		if node == nil {
			return nil
		}

		hostTotal, _ := tx.GostClientHost.Where(tx.GostClientHost.NodeCode.Eq(node.Code)).Count()
		if hostTotal > 0 {
			return errors.New("请先删除该节点的所有隧道")
		}

		forwardTotal, _ := tx.GostClientForward.Where(tx.GostClientForward.NodeCode.Eq(node.Code)).Count()
		if forwardTotal > 0 {
			return errors.New("请先删除该节点的所有隧道")
		}

		tunnelTotal, _ := tx.GostClientTunnel.Where(tx.GostClientTunnel.NodeCode.Eq(node.Code)).Count()
		if tunnelTotal > 0 {
			return errors.New("请先删除该节点的所有隧道")
		}

		if _, err := tx.GostNode.Where(tx.GostNode.Code.Eq(node.Code)).Delete(); err != nil {
			log.Error("删除节点失败", zap.Error(err))
			return errors.New("操作失败")
		}

		_, _ = tx.GostNodeConfig.Where(tx.GostNodeConfig.NodeCode.Eq(node.Code)).Delete()
		_, _ = tx.GostNodeLogger.Where(tx.GostNodeLogger.NodeCode.Eq(node.Code)).Delete()
		_, _ = tx.GostNodeBind.Where(tx.GostNodeBind.NodeCode.Eq(node.Code)).Delete()
		gost_engine.NodeStop(node.Code, "节点已被删除")
		cache.DelNodeInfo(node.Code)
		return nil
	})
}
