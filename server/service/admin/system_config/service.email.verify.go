package service

import (
	"errors"
	"server/model"
	"server/pkg/email"
	"server/pkg/utils"
	"server/repository/cache"
)

type EmailVerifyReq struct {
	Email string `binding:"required" json:"email"`
}

func (service *service) EmailVerify(req EmailVerifyReq) error {
	var cfg model.SystemConfigEmail
	cache.GetSystemConfigEmail(&cfg)
	if cfg.Enable != "1" {
		return errors.New("未启用服务")
	}

	return email.Send(email.Config{
		Host: cfg.Host,
		Port: utils.StrMustInt(cfg.Port),
		User: cfg.User,
		Pwd:  cfg.Pwd,
	}, email.Body{
		Mails:    []string{req.Email},
		NickName: cfg.NickName,
		Title:    "测试邮件",
		Body:     "test success",
	})
}
