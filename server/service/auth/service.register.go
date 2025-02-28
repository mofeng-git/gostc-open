package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/model"
	"server/pkg/utils"
	"server/repository"
	"server/service/common/cache"
)

type RegisterReq struct {
	Account      string `binding:"required" json:"account" label:"账号"`
	Password     string `binding:"required" json:"password" label:"秘密"`
	CaptchaKey   string `json:"captchaKey" label:"验证码Key"`
	CaptchaValue string `json:"captchaValue" label:"验证码Value"`
}

func (service *service) Register(ip string, req RegisterReq) (err error) {
	db, _, log := repository.Get("")
	if !cache.GetIpSecurity(ip) && !cache.ValidCaptcha(req.CaptchaKey, req.CaptchaValue, true) {
		return errors.New("验证码错误")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if tx.Where("account = ?", req.Account).First(&model.SystemUser{}).RowsAffected == 1 {
			return errors.New("该账号已被注册")
		}

		salt := utils.RandStr(8, utils.AllDict)
		if err = tx.Create(&model.SystemUser{
			Account:  req.Account,
			Password: utils.MD5AndSalt(req.Password, salt),
			Salt:     salt,
		}).Error; err != nil {
			log.Error("注册账号失败", zap.Error(err))
			return errors.New("注册失败，请联系管理员")
		}
		return nil
	})
}
