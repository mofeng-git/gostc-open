package service

import (
	"errors"
	"go.uber.org/zap"
	"server/pkg/jwt"
	"server/repository"
	"server/repository/cache"
	"server/repository/query"
)

type AdmissionReq struct {
	Code        string   `binding:"required" json:"code"`
	WhiteEnable int      `json:"whiteEnable"`
	White       []string `json:"white"`
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

		forward.WhiteEnable = req.WhiteEnable
		forward.SetWhiteList(req.White)

		if err := tx.GostClientForward.Save(forward); err != nil {
			log.Error("修改端口转发失败", zap.Error(err))
			return errors.New("操作失败")
		}
		cache.SetAdmissionInfo(cache.AdmissionInfo{
			Code:        forward.Code,
			WhiteEnable: req.WhiteEnable,
			WhiteList:   req.White,
		})
		return nil
	})
}
