package service

import (
	"errors"
	"go.uber.org/zap"
	"server/global"
	"server/model"
	"server/pkg/email"
	"server/pkg/utils"
	"server/repository"
	"server/repository/query"
	"server/service/common/cache"
	"time"
)

type ResetPwdReq struct {
	Account string `binding:"required" json:"account"`
	Key     string `binding:"required" json:"key"`
	Code    string `binding:"required" json:"code"`
}

func (service *service) ResetPwd(req ResetPwdReq) error {
	var cfg model.SystemConfigEmail
	cache.GetSystemConfigEmail(&cfg)
	if cfg.Enable != "1" {
		return errors.New("管理员未启用邮件服务")
	}

	code := cache.GetResetPwdEmailCode(req.Key+req.Account, false)
	if code != req.Code {
		return errors.New("验证码错误")
	}
	cache.GetResetPwdEmailCode(req.Key+req.Account, true)

	db, _, _ := repository.Get("")
	user, err := db.SystemUser.Preload(db.SystemUser.BindEmail).Where(db.SystemUser.Account.Eq(req.Account)).First()
	if err != nil {
		return errors.New("查询账号失败")
	}
	if user.BindEmail.Email == "" {
		return errors.New("该账号未绑定邮箱")
	}
	newPwd := utils.RandStr(8, utils.AllDict)
	user.Salt = utils.RandStr(8, utils.AllDict)
	user.Password = utils.MD5AndSalt(newPwd, user.Salt)

	return db.Transaction(func(tx *query.Query) error {
		if err := tx.SystemUser.Save(user); err != nil {
			global.Logger.Error("邮件重置验证码失败", zap.Error(err))
			return errors.New("重置失败")
		}

		if err := email.Send(email.Config{
			Host: cfg.Host,
			Port: utils.StrMustInt(cfg.Port),
			User: cfg.User,
			Pwd:  cfg.Pwd,
		}, email.Body{
			Mails:    []string{user.BindEmail.Email},
			NickName: cfg.NickName,
			Title:    "重置密码",
			Body:     cfg.GenerateResetPwdResultTpl(user.Account, newPwd, time.Now()),
		}); err != nil {
			global.Logger.Error("邮件重置密码，发送邮件失败")
			return errors.New("重置失败")
		}
		return nil
	})
}
