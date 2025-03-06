package service

import (
	"errors"
	"go.uber.org/zap"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/repository/query"
)

type PasswordReq struct {
	NewPwd string `binding:"required" json:"newPwd"`
	OldPwd string `binding:"required" json:"oldPwd"`
}

func (service *service) Password(claims jwt.Claims, req PasswordReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		user, err := tx.SystemUser.Where(tx.SystemUser.Code.Eq(claims.Code)).First()
		if err != nil {
			return errors.New("账号错误")
		}
		if user.Password != utils.MD5AndSalt(req.OldPwd, user.Salt) {
			return errors.New("原密码错误")
		}
		user.Password = utils.MD5AndSalt(req.NewPwd, user.Salt)
		if err := tx.SystemUser.Save(user); err != nil {
			log.Error("修改密码失败", zap.Error(err))
			return errors.New("修改失败")
		}
		return nil
	})
}
