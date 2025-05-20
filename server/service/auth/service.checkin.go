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
	"server/service/common/cache"
	"time"
)

func (service *service) Checkin(claims jwt.Claims) (err error) {
	var cfg model.SystemConfigBase
	cache.GetSystemConfigBase(&cfg)
	if cfg.CheckIn != "1" {
		return errors.New("未启用签到功能")
	}

	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		user, err := tx.SystemUser.Where(tx.SystemUser.Code.Eq(claims.Code)).First()
		if err != nil {
			return errors.New("账号错误")
		}

		checkin, _ := tx.SystemUserCheckin.Where(
			tx.SystemUserCheckin.UserCode.Eq(claims.Code),
			tx.SystemUserCheckin.EventDate.Eq(time.Now().Format(time.DateOnly)),
		).First()
		if checkin != nil {
			return errors.New("已签到")
		}

		amount := decimal.NewFromInt(int64(utils.RandNum(cfg.CheckInEnd-cfg.CheckInStart) + cfg.CheckInStart))
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
