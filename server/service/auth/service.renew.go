package service

import (
	"errors"
	"go.uber.org/zap"
	"server/global"
	"server/pkg/jwt"
	"server/repository"
	"strconv"
	"time"
)

func (service *service) Renew(oldToken string, claims jwt.Claims) (result LoginResp, err error) {
	db, _, _ := repository.Get("")
	var bufExpAt int64 = 60 * 60 * 2
	if claims.ExpiresAt < time.Now().Unix() {
		return result, errors.New("登录失效")
	}
	if claims.ExpiresAt-bufExpAt < time.Now().Unix() {
		return LoginResp{
			Token: oldToken,
			ExpAt: claims.ExpiresAt,
		}, nil
	}

	user, _ := db.SystemUser.Where(db.SystemUser.Code.Eq(claims.Code)).First()
	if user == nil {
		return result, errors.New("未查询到账户信息")
	}

	token, err := global.Jwt.GenerateToken(global.Jwt.NewClaims(user.Code, map[string]string{
		"admin": strconv.Itoa(user.Admin),
	}, global.Config.AuthExp))
	if err != nil {
		global.Logger.Error("生成Token失败", zap.Error(err))
		return LoginResp{}, errors.New("登录失败，请联系管理员")
	}

	return LoginResp{
		Token: token,
		ExpAt: time.Now().Add(global.Config.AuthExp).Unix(),
	}, nil
}
