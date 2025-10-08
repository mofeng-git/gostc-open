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
	Code           string `binding:"required" json:"code"`
	Name           string `binding:"required" json:"name"`
	TargetIp       string `binding:"required" json:"targetIp"`
	TargetPort     string `binding:"required" json:"targetPort"`
	Forward        int    `json:"forward"`
	UseEncryption  int    `json:"useEncryption"`
	UseCompression int    `json:"useCompression"`
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

		p2p, _ := tx.GostClientP2P.Where(tx.GostClientP2P.Code.Eq(req.Code), tx.GostClientP2P.UserCode.Eq(claims.Code)).First()
		if p2p == nil {
			return errors.New("操作失败")
		}

		p2p.Name = req.Name
		p2p.TargetIp = req.TargetIp
		p2p.TargetPort = req.TargetPort
		p2p.Forward = req.Forward
		p2p.UseEncryption = req.UseEncryption
		p2p.UseCompression = req.UseCompression

		if err := tx.GostClientP2P.Save(p2p); err != nil {
			log.Error("修改P2P隧道失败", zap.Error(err))
			return errors.New("操作失败")
		}
		engine.ClientP2PConfig(tx, p2p.Code)
		return nil
	})
}
