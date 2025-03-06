package service

import (
	"errors"
	"go.uber.org/zap"
	"server/pkg/jwt"
	"server/repository"
	"server/repository/query"
)

func (service *service) CloseOtp(claims jwt.Claims) (err error) {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		user, err := tx.SystemUser.Where(tx.SystemUser.Code.Eq(claims.Code)).First()
		if err != nil {
			return errors.New("账号错误")
		}
		if user.OtpKey == "" {
			return nil
		}
		user.OtpKey = ""
		if err := tx.SystemUser.Save(user); err != nil {
			log.Error("关闭otp失败", zap.Error(err))
			return errors.New("操作失败")
		}
		return nil
	})
}
