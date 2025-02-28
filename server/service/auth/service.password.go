package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/model"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
)

type PasswordReq struct {
	NewPwd string `binding:"required" json:"newPwd"`
	OldPwd string `binding:"required" json:"oldPwd"`
}

func (service *service) Password(claims jwt.Claims, req PasswordReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *gorm.DB) error {
		var user model.SystemUser
		if tx.Where("code = ?", claims.Code).First(&user).RowsAffected == 0 {
			return errors.New("用户错误")
		}
		if user.Password != utils.MD5AndSalt(req.OldPwd, user.Salt) {
			return errors.New("原密码错误")
		}
		user.Password = utils.MD5AndSalt(req.NewPwd, user.Salt)
		if err := tx.Save(&user).Error; err != nil {
			log.Error("修改密码失败", zap.Error(err))
			return errors.New("修改失败")
		}
		return nil
	})
}
