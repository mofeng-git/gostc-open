package service

import (
	"errors"
	"go.uber.org/zap"
	"server/global"
	"server/model"
	"server/pkg/jwt"
	"server/repository"
	cache2 "server/repository/cache"
	"server/repository/query"
)

type BindEmailReq struct {
	Target string `binding:"required" json:"target"`
	Code   string `binding:"required" json:"code"`
	Key    string `binding:"required" json:"key"`
}

func (service *service) BindEmail(claims jwt.Claims, req BindEmailReq) error {
	var cfg model.SystemConfigEmail
	cache2.GetSystemConfigEmail(&cfg)
	if cfg.Enable != "1" {
		return errors.New("管理员未启用邮件服务")
	}

	code := cache2.GetBindEmailCode(req.Key, false)
	if code != req.Code {
		return errors.New("验证码错误")
	}
	cache2.GetBindEmailCode(req.Key, true)

	db, _, _ := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		{
			bindEmail, _ := tx.SystemUserEmail.Where(tx.SystemUserEmail.Email.Eq(req.Target)).First()
			if bindEmail != nil {
				return errors.New("该邮箱已被其他账号绑定")
			}
		}
		{
			bindEmail, _ := tx.SystemUserEmail.Where(tx.SystemUserEmail.UserCode.Eq(claims.Code)).First()
			if bindEmail != nil {
				return errors.New("当前账号已绑定邮箱")
			}
		}
		if err := tx.SystemUserEmail.Create(&model.SystemUserEmail{
			Email:    req.Target,
			UserCode: claims.Code,
		}); err != nil {
			global.Logger.Error("绑定邮箱失败", zap.Error(err))
			return errors.New("绑定失败")
		}
		return nil
	})
}
