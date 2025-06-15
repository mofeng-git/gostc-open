package service

import (
	"errors"
	"fmt"
	"server/model"
	"server/pkg/email"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	cache2 "server/repository/cache"
	"server/repository/query"
	"time"
)

type GenBindEmailCodeReq struct {
	Target string `binding:"required" json:"target"`
}

func (service *service) GenBindEmailCode(claims jwt.Claims, req GenBindEmailCodeReq) (string, error) {
	var cfg model.SystemConfigEmail
	cache2.GetSystemConfigEmail(&cfg)
	if cfg.Enable != "1" {
		return "", errors.New("管理员未启用邮件服务")
	}

	lastTime := cache2.GetBindEmailLastTime(claims.Code)
	space := lastTime.Add(time.Minute).Sub(time.Now()).Milliseconds()
	if space > 0 {
		return "", fmt.Errorf("请等待%d秒后再尝试", space/1000)
	}

	db, _, _ := repository.Get("")
	if err := db.Transaction(func(tx *query.Query) error {
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
		return nil
	}); err != nil {
		return "", err
	}

	str := utils.RandStr(32, utils.AllDict)
	code := utils.RandStr(6, utils.NumDict)
	cache2.SetBindEmailCode(str, code, time.Minute*5)
	if err := email.Send(email.Config{
		Host: cfg.Host,
		Port: utils.StrMustInt(cfg.Port),
		User: cfg.User,
		Pwd:  cfg.Pwd,
	}, email.Body{
		Mails:    []string{req.Target},
		NickName: cfg.NickName,
		Title:    "绑定邮箱",
		Body:     cfg.GenerateUserBindEmailTpl(code, time.Now()),
	}); err != nil {
		return "", err
	}
	cache2.SetBindEmailLastTime(claims.Code)
	return str, nil
}
