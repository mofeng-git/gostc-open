package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
	"image/png"
	"server/model"
	"server/global"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/service/common/cache"
	"time"
)

type GenOtpResp struct {
	Key string `json:"key"`
	Img string `json:"img"`
}

func (service *service) GenOtp(claims jwt.Claims) (result GenOtpResp, err error) {
	db, _, _ := repository.Get("")
	var user model.SystemUser
	if db.Where("code = ?", claims.Code).First(&user).RowsAffected == 0 {
		return result, errors.New("用户不存在")
	}
	if user.OtpKey != "" {
		return result, errors.New("已开启二步验证")
	}

	var list []model.SystemConfig
	db.Where("kind = ?", model.SYSTEM_CONFIG_KIND_BASE).Find(&list)
	systemConfig := model.GetSystemConfigBase(list)

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
	cache.SetBindOtp(result.Key, optKey.Secret(), time.Minute*5)
	return result, nil
}
