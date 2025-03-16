package service

import (
	"errors"
	"go.uber.org/zap"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/repository/query"
	"server/service/gost_engine"
)

type UpdateReq struct {
	Code     string `binding:"required" json:"code"`
	Name     string `binding:"required" json:"name"`
	Port     string `binding:"required" json:"port" label:"本地端口"`
	Protocol string `binding:"required" json:"protocol" label:"协议"`
}

func (service *service) Update(claims jwt.Claims, req UpdateReq) error {
	db, _, log := repository.Get("")
	if !utils.ValidatePort(req.Port) {
		return errors.New("内网端口格式错误")
	}

	return db.Transaction(func(tx *query.Query) error {
		user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(claims.Code)).First()
		if user == nil {
			return errors.New("用户错误")
		}

		proxy, _ := tx.GostClientProxy.Where(
			tx.GostClientProxy.UserCode.Eq(user.Code),
			tx.GostClientProxy.Code.Eq(req.Code),
		).First()
		if proxy == nil {
			return errors.New("操作失败")
		}

		client, _ := tx.GostClient.Where(
			tx.GostClient.UserCode.Eq(claims.Code),
			tx.GostClient.Code.Eq(proxy.ClientCode),
		).First()
		if client == nil {
			return errors.New("客户端错误")
		}

		if proxy.Port != req.Port {
			if err := CheckPort(tx, *client, req.Port); err != nil {
				return err
			}
		}

		proxy.Name = req.Name
		proxy.Port = req.Port
		proxy.Protocol = req.Protocol

		if err := tx.GostClientProxy.Save(proxy); err != nil {
			log.Error("修改代理隧道失败", zap.Error(err))
			return errors.New("操作失败")
		}
		gost_engine.ClientProxyConfig(tx, proxy.Code)
		return nil
	})
}
