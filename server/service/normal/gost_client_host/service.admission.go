package service

import (
	"errors"
	"go.uber.org/zap"
	"server/pkg/jwt"
	"server/repository"
	"server/repository/query"
	"server/service/gost_engine"
)

type AdmissionReq struct {
	Code        string   `binding:"required" json:"code"`
	WhiteEnable int      `json:"whiteEnable"`
	BlackEnable int      `json:"blackEnable"`
	White       []string `json:"white"`
	Black       []string `json:"black"`
}

func (service *service) Admission(claims jwt.Claims, req AdmissionReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(claims.Code)).First()
		if user == nil {
			return errors.New("用户错误")
		}

		host, _ := tx.GostClientHost.Where(
			tx.GostClientHost.UserCode.Eq(user.Code),
			tx.GostClientHost.Code.Eq(req.Code),
		).First()
		if host == nil {
			return errors.New("操作失败")
		}

		host.WhiteEnable = req.WhiteEnable
		host.BlackEnable = req.BlackEnable
		host.SetWhiteList(req.White)
		host.SetBlackList(req.Black)

		if err := tx.GostClientHost.Save(host); err != nil {
			log.Error("修改域名解析失败", zap.Error(err))
			return errors.New("操作失败")
		}
		gost_engine.ClientHostConfig(tx, host.Code)
		return nil
	})
}
