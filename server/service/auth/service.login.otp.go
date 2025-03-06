package service

import (
	"errors"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
	"server/global"
	"server/repository"
	"server/service/common/cache"
	"strconv"
	"time"
)

type LoginOtpReq struct {
	Key   string `binding:"required" json:"key"`
	Value string `binding:"required" json:"value"`
}

type LoginOtpResp struct {
	Token string `json:"token"`
	ExpAt int64  `json:"expAt"`
}

func (service *service) LoginOtp(ip string, req LoginOtpReq) (result LoginOtpResp, err error) {
	db, _, _ := repository.Get("")
	userCode := cache.GetLoginOtp(req.Key, true)
	if userCode == "" {
		return result, errors.New("登录失败")
	}
	user, _ := db.SystemUser.Where(db.SystemUser.Code.Eq(userCode)).First()
	if user == nil {
		return result, errors.New("未查询到账户信息")
	}
	// 进行二步验证
	if ok, _ := totp.ValidateCustom(req.Value, user.OtpKey, time.Now(), totp.ValidateOpts{
		Period: 30,
		Digits: otp.DigitsSix,
	}); !ok {
		return result, errors.New("验证失败")
	}

	token, err := global.Jwt.GenerateToken(global.Jwt.NewClaims(user.Code, map[string]string{
		"admin": strconv.Itoa(user.Admin),
	}, global.Config.AuthExp))
	if err != nil {
		global.Logger.Error("生成Token失败", zap.Error(err))
		return result, errors.New("登录失败，请联系管理员")
	}
	cache.SetIpSecurity(ip, true)
	return LoginOtpResp{
		Token: token,
		ExpAt: time.Now().Add(global.Config.AuthExp).Unix(),
	}, nil
}
