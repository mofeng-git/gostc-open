package service

import (
	"errors"
	"go.uber.org/zap"
	"server/pkg/jwt"
	"server/repository"
	"server/repository/query"
	"server/service/engine"
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

		forward, _ := tx.GostClientForward.Where(
			tx.GostClientForward.UserCode.Eq(user.Code),
			tx.GostClientForward.Code.Eq(req.Code),
		).First()
		if forward == nil {
			return errors.New("操作失败")
		}

		//forward.WhiteEnable = req.WhiteEnable
		//forward.BlackEnable = req.BlackEnable
		//forward.SetWhiteList(req.White)
		//forward.SetBlackList(req.Black)

		if err := tx.GostClientForward.Save(forward); err != nil {
			log.Error("修改端口转发失败", zap.Error(err))
			return errors.New("操作失败")
		}
		engine.ClientForwardConfig(tx, forward.Code)
		return nil
	})
}
