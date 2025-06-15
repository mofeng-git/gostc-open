package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
	"image/png"
	"server/global"
	"server/model"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	cache2 "server/repository/cache"
	"time"
)

type GenOtpResp struct {
	Key string `json:"key"`
	Img string `json:"img"`
}

func (service *service) GenOtp(claims jwt.Claims) (result GenOtpResp, err error) {
	db, _, _ := repository.Get("")
	user, err := db.SystemUser.Where(db.SystemUser.Code.Eq(claims.Code)).First()
	if err != nil {
		return result, errors.New("账号错误")
	}
	if user.OtpKey != "" {
		return result, errors.New("已开启二步验证")
	}

	var systemConfig model.SystemConfigBase
	cache2.GetSystemConfigBase(&systemConfig)

	result.Key = utils.RandStr(32, utils.AllDict)
	optKey, err := totp.Generate(totp.GenerateOpts{
		Issuer:      systemConfig.Title,
		AccountName: user.Account,
	})
	if err != nil {
		global.Logger.Error("生成二步验证密钥失败", zap.Error(err))
		return result, errors.New("生成二维码失败")
	}
	var imgBuf bytes.Buffer
	optImg, err := optKey.Image(200, 200)
	if err != nil {
		return result, errors.New("生成二维码失败")
	}
	if err = png.Encode(&imgBuf, optImg); err != nil {
		return result, errors.New("生成二维码失败")
	}
	result.Img = base64.StdEncoding.EncodeToString(imgBuf.Bytes())
	cache2.SetBindOtp(result.Key, optKey.Secret(), time.Minute*5)
	return result, nil
}
