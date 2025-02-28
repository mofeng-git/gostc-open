package service

import (
	"errors"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/model"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"time"
)

func (service *service) Checkin(claims jwt.Claims) (err error) {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *gorm.DB) error {
		var user model.SystemUser
		if tx.Where("code = ?", claims.Code).First(&user).RowsAffected == 0 {
			return errors.New("账号错误")
		}
		if tx.Where("user_code = ? AND event_date = ?", user.Code, time.Now().Format(time.DateOnly)).First(&model.SystemUserCheckin{}).RowsAffected != 0 {
			return errors.New("已签到")
		}
		amount := decimal.NewFromInt(int64(utils.RandNum(5) + 1))
		user.Amount = user.Amount.Add(amount)
		if err := tx.Save(&user).Error; err != nil {
			log.Error("签到失败", zap.Error(err))
			return errors.New("签到失败")
		}
		if err := tx.Create(&model.SystemUserCheckin{
			UserCode:  user.Code,
			Account:   user.Account,
			EventDate: time.Now().Format(time.DateOnly),
			Amount:    amount,
		}).Error; err != nil {
			log.Error("签到失败", zap.Error(err))
			return errors.New("签到失败")
		}
		return nil
	})
}
