package service

import (
	"errors"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/model"
	"server/pkg/jwt"
	"server/repository"
	"server/service/common/cache"
)

type OpenOtpReq struct {
	Key   string `binding:"required" json:"key"`
	Value string `binding:"required" json:"value"`
}

func (service *service) OpenOtp(claims jwt.Claims, req OpenOtpReq) (err error) {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *gorm.DB) error {
		var user model.SystemUser
		if tx.Where("code = ?", claims.Code).First(&user).RowsAffected == 0 {
			return errors.New("用户不存在")
		}
		if user.OtpKey != "" {
			return errors.New("已开启二步验证")
		}
		otpKey := cache.GetBindOtp(req.Key, true)
		if otpKey == "" {
			return errors.New("二维码已失效")
		}
		if !totp.Validate(req.Value, otpKey) {
			return errors.New("验证失败")
		}
		user.OtpKey = otpKey
		if err := tx.Save(&user).Error; err != nil {
			log.Error("开启otp失败", zap.Error(err))
			return errors.New("操作失败")
		}
		return nil
	})
}
