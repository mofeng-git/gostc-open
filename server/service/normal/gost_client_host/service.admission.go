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

		host, _ := tx.GostClientHost.Where(
			tx.GostClientHost.UserCode.Eq(user.Code),
			tx.GostClientHost.Code.Eq(req.Code),
		).First()
		if host == nil {
			return errors.New("操作失败")
		}

		host.WhiteEnable = req.WhiteEnable
		host.SetWhiteList(req.White)

		if err := tx.GostClientHost.Save(host); err != nil {
			log.Error("修改域名解析失败", zap.Error(err))
			return errors.New("操作失败")
		}
		cache.SetAdmissionInfo(cache.AdmissionInfo{
			Code:        host.Code,
			WhiteEnable: req.WhiteEnable,
			WhiteList:   req.White,
		})
		return nil
	})
}
