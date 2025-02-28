package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/model"
	"server/pkg/jwt"
	"server/repository"
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
	return db.Transaction(func(tx *gorm.DB) error {
		var user model.SystemUser
		if tx.Where("code = ?", claims.Code).First(&user).RowsAffected == 0 {
			return errors.New("用户错误")
		}

		var forward model.GostClientForward
		if tx.Where("code = ? AND user_code = ?", req.Code, user.Code).First(&forward).RowsAffected == 0 {
			return errors.New("操作失败")
		}

		forward.WhiteEnable = req.WhiteEnable
		forward.BlackEnable = req.BlackEnable
		forward.SetWhiteList(req.White)
		forward.SetBlackList(req.Black)

		if err := tx.Save(&forward).Error; err != nil {
			log.Error("修改端口转发失败", zap.Error(err))
			return errors.New("操作失败")
		}
		gost_engine.ClientForwardConfig(tx, forward.Code)
		return nil
	})
}
