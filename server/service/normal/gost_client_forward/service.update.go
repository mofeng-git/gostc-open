package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/model"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/service/gost_engine"
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

	return db.Transaction(func(tx *gorm.DB) error {
		var user model.SystemUser
		if tx.Where("code = ?", claims.Code).First(&user).RowsAffected == 0 {
			return errors.New("用户错误")
		}

		var forward model.GostClientForward
		if tx.Where("code = ? AND user_code = ?", req.Code, user.Code).First(&forward).RowsAffected == 0 {
			return errors.New("操作失败")
		}

		forward.Name = req.Name
		forward.TargetIp = req.TargetIp
		forward.TargetPort = req.TargetPort
		forward.NoDelay = req.NoDelay

		if err := tx.Save(&forward).Error; err != nil {
			log.Error("修改端口转发失败", zap.Error(err))
			return errors.New("操作失败")
		}
		gost_engine.ClientForwardConfig(tx, forward.Code)
		return nil
	})
}
