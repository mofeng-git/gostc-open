package service

import (
	"errors"
	"go.uber.org/zap"
	"server/repository"
	"server/repository/cache"
	"server/repository/query"
	forwardService "server/service/admin/gost_client_forward"
	hostService "server/service/admin/gost_client_host"
	p2pService "server/service/admin/gost_client_p2p"
	proxyService "server/service/admin/gost_client_proxy"
	tunnelService "server/service/admin/gost_client_tunnel"
	"server/service/engine"
)

type DeleteReq struct {
	Code string `binding:"required" json:"code" label:"编号"`
}

func (service *service) Delete(req DeleteReq) error {
	db, _, log := repository.Get("")
	if err := db.Transaction(func(tx *query.Query) error {
		node, _ := tx.GostNode.Where(tx.GostNode.Code.Eq(req.Code)).First()
		if node == nil {
			return nil
		}

		if _, err := tx.GostNode.Where(tx.GostNode.Code.Eq(node.Code)).Delete(); err != nil {
			log.Error("删除节点失败", zap.Error(err))
			return errors.New("操作失败")
		}

		_, _ = tx.GostNodeConfig.Where(tx.GostNodeConfig.NodeCode.Eq(node.Code)).Delete()
		_, _ = tx.GostNodeBind.Where(tx.GostNodeBind.NodeCode.Eq(node.Code)).Delete()
		engine.NodeStop(node.Code, "节点已被删除")
		cache.DelNodeInfo(node.Code)
		return nil
	}); err != nil {
		return err
	}

	var hostCodes []string
	_ = db.GostClientHost.Where(db.GostClientHost.NodeCode.Eq(req.Code)).Pluck(db.GostClientHost.Code, &hostCodes)
	for _, code := range hostCodes {
		_ = hostService.Service.Delete(hostService.DeleteReq{Code: code})
	}
	var forwardCodes []string
	_ = db.GostClientForward.Where(db.GostClientForward.NodeCode.Eq(req.Code)).Pluck(db.GostClientForward.Code, &forwardCodes)
	for _, code := range forwardCodes {
		_ = forwardService.Service.Delete(forwardService.DeleteReq{Code: code})
	}
	var tunnelCodes []string
	_ = db.GostClientTunnel.Where(db.GostClientTunnel.NodeCode.Eq(req.Code)).Pluck(db.GostClientTunnel.Code, &tunnelCodes)
	for _, code := range tunnelCodes {
		_ = tunnelService.Service.Delete(tunnelService.DeleteReq{Code: code})
	}
	var p2pCodes []string
	_ = db.GostClientP2P.Where(db.GostClientP2P.NodeCode.Eq(req.Code)).Pluck(db.GostClientP2P.Code, &p2pCodes)
	for _, code := range p2pCodes {
		_ = p2pService.Service.Delete(p2pService.DeleteReq{Code: code})
	}
	var proxyCodes []string
	_ = db.GostClientProxy.Where(db.GostClientProxy.NodeCode.Eq(req.Code)).Pluck(db.GostClientProxy.Code, &proxyCodes)
	for _, code := range proxyCodes {
		_ = proxyService.Service.Delete(proxyService.DeleteReq{Code: code})
	}
	return nil
}
