package service

import (
	"errors"
	"go.uber.org/zap"
	"server/pkg/jwt"
	"server/repository"
	"server/repository/query"
	"server/service/engine"
)

type DeleteReq struct {
	Code string `binding:"required" json:"code" label:"编号"`
}

func (service *service) Delete(claims jwt.Claims, req DeleteReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		client, err := tx.GostClient.Where(tx.GostClient.Code.Eq(req.Code), tx.GostClient.UserCode.Eq(claims.Code)).First()
		if err != nil {
			return nil
		}
		hostTotal, _ := tx.GostClientHost.Where(tx.GostClientHost.ClientCode.Eq(client.Code)).Count()
		if hostTotal > 0 {
			return errors.New("请先删除该客户端的所有隧道")
		}

		forwardTotal, _ := tx.GostClientForward.Where(tx.GostClientForward.ClientCode.Eq(client.Code)).Count()
		if forwardTotal > 0 {
			return errors.New("请先删除该客户端的所有隧道")
		}

		tunnelTotal, _ := tx.GostClientTunnel.Where(tx.GostClientTunnel.ClientCode.Eq(client.Code)).Count()
		if tunnelTotal > 0 {
			return errors.New("请先删除该客户端的所有隧道")
		}

		p2pTotal, _ := tx.GostClientP2P.Where(tx.GostClientP2P.ClientCode.Eq(client.Code)).Count()
		if p2pTotal > 0 {
			return errors.New("请先删除该客户端的所有隧道")
		}

		proxyTotal, _ := tx.GostClientProxy.Where(tx.GostClientProxy.ClientCode.Eq(client.Code)).Count()
		if proxyTotal > 0 {
			return errors.New("请先删除该客户端的所有隧道")
		}

		cfgTotal, _ := tx.FrpClientCfg.Where(tx.FrpClientCfg.ClientCode.Eq(client.Code)).Count()
		if cfgTotal > 0 {
			return errors.New("请先删除该客户端的所有隧道")
		}

		if _, err := tx.GostClient.Where(tx.GostClient.Code.Eq(client.Code)).Delete(); err != nil {
			log.Error("删除客户端失败", zap.Error(err))
			return errors.New("操作失败")
		}
		engine.ClientStop(client.Code, "客户端已被删除")
		return nil
	})
}
