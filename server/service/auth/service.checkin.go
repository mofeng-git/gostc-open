package service

import (
	"errors"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"server/model"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/repository/query"
	"time"
)

func (service *service) Checkin(claims jwt.Claims) (err error) {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		user, err := tx.SystemUser.Where(tx.SystemUser.Code.Eq(claims.Code)).First()
		if err != nil {
			return errors.New("账号错误")
		}

		checkin, _ := tx.SystemUserCheckin.Where(
			tx.SystemUserCheckin.UserCode.Eq(claims.Code),
			tx.SystemUserCheckin.EventDate.Eq(time.Now().Format(time.DateTime)),
		).First()
		if checkin != nil {
			return errors.New("已签到")
		}

		amount := decimal.NewFromInt(int64(utils.RandNum(5) + 1))
		user.Amount = user.Amount.Add(amount)
		if err := tx.SystemUser.Save(user); err != nil {
			log.Error("签到失败", zap.Error(err))
			return errors.New("签到失败")
		}
		if err := tx.SystemUserCheckin.Create(&model.SystemUserCheckin{
			UserCode:  user.Code,
			Account:   user.Account,
			EventDate: time.Now().Format(time.DateOnly),
			Amount:    amount,
		}); err != nil {
			log.Error("签到失败", zap.Error(err))
			return errors.New("签到失败")
		}
		return nil
	})
}
