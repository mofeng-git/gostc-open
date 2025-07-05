package service

import (
	"errors"
	"go.uber.org/zap"
	"server/global"
	"server/pkg/utils"
	"server/repository"
	cache2 "server/repository/cache"
	"strconv"
	"time"
)

type LoginReq struct {
	Account      string `binding:"required" json:"account" label:"账号"`
	Password     string `binding:"required" json:"password" label:"秘密"`
	CaptchaKey   string `json:"captchaKey" label:"验证码Key"`
	CaptchaValue string `json:"captchaValue" label:"验证码Value"`
}

type LoginResp struct {
	Otp   int    `json:"otp"`
	Token string `json:"token"`
	ExpAt int64  `json:"expAt"`
}

func (service *service) Login(ip string, req LoginReq) (result LoginResp, err error) {
	defer func() {
		if err != nil {
			cache2.SetIpSecurity(ip, false)
		} else {
			cache2.SetIpSecurity(ip, true)
		}
	}()
	db, _, _ := repository.Get("")
	if !cache2.GetIpSecurity(ip) && !cache2.ValidCaptcha(req.CaptchaKey, req.CaptchaValue, true) {
		return result, errors.New("验证码错误")
	}

	user, _ := db.SystemUser.Where(db.SystemUser.Account.Eq(req.Account)).First()
	if user == nil {
		return result, errors.New("未查询到账户信息")
	}

	if utils.MD5AndSalt(req.Password, user.Salt) != user.Password {
		return result, errors.New("账号或密码错误")
	}

	if user.OtpKey == "" {
		token, err := global.Jwt.GenerateToken(global.Jwt.NewClaims(user.Code, map[string]string{
			"admin": strconv.Itoa(user.Admin),
		}, global.Config.AuthExp))
		if err != nil {
			global.Logger.Error("生成Token失败", zap.Error(err))
			return LoginResp{}, errors.New("登录失败，请联系管理员")
		}
		result = LoginResp{
			Otp:   2,
			Token: token,
			ExpAt: time.Now().Add(global.Config.AuthExp).Unix(),
		}
	} else {
		key := utils.RandStr(32, utils.AllDict)
		cache2.SetLoginOtp(key, user.Code, time.Minute*5)
		result = LoginResp{
			Otp:   1,
			Token: key,
			ExpAt: 0,
		}
	}
	return result, nil
}
