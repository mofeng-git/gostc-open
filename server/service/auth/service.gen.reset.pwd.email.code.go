package service

import (
	"errors"
	"fmt"
	"server/model"
	"server/pkg/email"
	"server/pkg/utils"
	"server/repository"
	"server/service/common/cache"
	"time"
)

type GenResetPwdEmailCodeReq struct {
	Account      string `binding:"required" json:"account"`
	CaptchaKey   string `binding:"required" json:"captchaKey" label:"验证码Key"`
	CaptchaValue string `binding:"required" json:"captchaValue" label:"验证码Value"`
}

func (service *service) GenResetPwdEmailCode(req GenResetPwdEmailCodeReq) (string, error) {
	var cfg model.SystemConfigEmail
	cache.GetSystemConfigEmail(&cfg)
	if cfg.Enable != "1" {
		return "", errors.New("管理员未启用邮件服务")
	}

	if !cache.ValidCaptcha(req.CaptchaKey, req.CaptchaValue, true) {
		return "", errors.New("验证码错误")
	}

	lastTime := cache.GetResetPwdEmailLastTime(req.Account)
	space := lastTime.Add(time.Minute).Sub(time.Now()).Milliseconds()
	if space > 0 {
		return "", fmt.Errorf("请等待%d秒后再尝试", space/1000)
	}

	db, _, _ := repository.Get("")
	user, err := db.SystemUser.Preload(db.SystemUser.BindEmail).Where(db.SystemUser.Account.Eq(req.Account)).First()
	if err != nil {
		return "", errors.New("查询账号失败")
	}
	if user.BindEmail.Email == "" {
		return "", errors.New("该账号未绑定邮箱")
	}

	str := utils.RandStr(32, utils.AllDict)
	code := utils.RandStr(6, utils.NumDict)
	cache.SetResetPwdEmailCode(str+req.Account, code, time.Minute*5)
	if err := email.Send(email.Config{
		Host: cfg.Host,
		Port: utils.StrMustInt(cfg.Port),
		User: cfg.User,
		Pwd:  cfg.Pwd,
	}, email.Body{
		Mails:    []string{user.BindEmail.Email},
		NickName: cfg.NickName,
		Title:    "重置密码",
		Body:     cfg.GenerateResetPwdTpl(code, time.Now()),
	}); err != nil {
		return "", err
	}
	cache.SetResetPwdEmailLastTime(req.Account)
	return str, nil
}
