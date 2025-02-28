package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/model"
	"server/pkg/jwt"
	"server/repository"
)

func (service *service) CloseOtp(claims jwt.Claims) (err error) {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *gorm.DB) error {
		var user model.SystemUser
		if tx.Where("code = ?", claims.Code).First(&user).RowsAffected == 0 {
			return errors.New("用户不存在")
		}
		if user.OtpKey == "" {
			return nil
		}
		user.OtpKey = ""
		if err := tx.Save(&user).Error; err != nil {
			log.Error("关闭otp失败", zap.Error(err))
			return errors.New("操作失败")
		}
		return nil
	})
}
