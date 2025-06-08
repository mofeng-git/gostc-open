package service

import (
	"errors"
	"go.uber.org/zap"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/repository/query"
	"server/service/engine"
)

type UpdateReq struct {
	Code       string `binding:"required" json:"code"`
	Name       string `binding:"required" json:"name"`
	TargetIp   string `binding:"required" json:"targetIp"`
	TargetPort string `binding:"required" json:"targetPort"`
	NoDelay    int    `json:"noDelay" label:"兼容模式"`
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

		tunnel, _ := tx.GostClientTunnel.Where(tx.GostClientTunnel.Code.Eq(req.Code), tx.GostClientTunnel.UserCode.Eq(claims.Code)).First()
		if tunnel == nil {
			return errors.New("操作失败")
		}

		tunnel.Name = req.Name
		tunnel.TargetIp = req.TargetIp
		tunnel.TargetPort = req.TargetPort
		tunnel.NoDelay = req.NoDelay

		if err := tx.GostClientTunnel.Save(tunnel); err != nil {
			log.Error("修改私有隧道失败", zap.Error(err))
			return errors.New("操作失败")
		}
		engine.ClientTunnelConfig(tx, tunnel.Code)
		return nil
	})
}
