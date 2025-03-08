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
	Code          string `binding:"required" json:"code"`
	Name          string `binding:"required" json:"name"`
	TargetIp      string `binding:"required" json:"targetIp"`
	TargetPort    string `binding:"required" json:"targetPort"`
	ProxyProtocol int    `json:"proxyProtocol"`
}

func (service *service) Update(claims jwt.Claims, req UpdateReq) error {
	db, _, log := repository.Get("")
	if !utils.ValidateLocalIP(req.TargetIp) {
		return errors.New("内网IP格式错误")
	}
	if !utils.ValidatePort(req.TargetPort) {
		return errors.New("内网端口格式错误")
	}

	return db.Transaction(func(tx *query.Query) error {
		user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(claims.Code)).First()
		if user == nil {
			return errors.New("用户错误")
		}

		forward, _ := tx.GostClientForward.Where(
			tx.GostClientForward.UserCode.Eq(user.Code),
			tx.GostClientForward.Code.Eq(req.Code),
		).First()
		if forward == nil {
			return errors.New("操作失败")
		}

		forward.Name = req.Name
		forward.TargetIp = req.TargetIp
		forward.TargetPort = req.TargetPort
		forward.ProxyProtocol = req.ProxyProtocol

		if err := tx.GostClientForward.Save(forward); err != nil {
			log.Error("修改端口转发失败", zap.Error(err))
			return errors.New("操作失败")
		}
		gost_engine.ClientForwardConfig(tx, forward.Code)
		return nil
	})
}
