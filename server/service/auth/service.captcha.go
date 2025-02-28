package service

import (
	"encoding/base64"
	"errors"
	"go.uber.org/zap"
	"server/pkg/captcha"
	"server/pkg/utils"
	"server/repository"
	"server/service/common/cache"
	"time"
)

type CaptchaResp struct {
	Key      string `json:"key"`
	Bs64     string `json:"bs64"`
	Security bool   `json:"security"`
}

func (service *service) Captcha(ip string) (result CaptchaResp, err error) {
	_, _, log := repository.Get("")
	code := utils.RandStr(4, utils.NumDict)
	bytes, err := captcha.Generate(120, 40, code)
	if err != nil {
		log.Error("生成图片验证码失败", zap.Error(err))
		return result, errors.New("获取失败")
	}
	result.Bs64 = base64.StdEncoding.EncodeToString(bytes)
	result.Key = utils.RandStr(16, utils.AllDict)
	cache.SetCaptcha(result.Key, code, time.Minute*5)
	result.Security = cache.GetIpSecurity(ip)
	return result, nil
}
