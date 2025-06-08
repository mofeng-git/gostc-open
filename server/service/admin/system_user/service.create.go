package service

import (
	"errors"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"server/model"
	"server/pkg/utils"
	"server/repository"
)

type CreateReq struct {
	Account  string          `binding:"required" json:"account" label:"账号"`
	Password string          `binding:"required" json:"password" label:"密码"`
	Amount   decimal.Decimal `binding:"required" json:"amount" label:"积分"`
	Admin    int             `binding:"required" json:"admin"`
}

func (service *service) Create(req CreateReq) error {
	db, _, log := repository.Get("")
	user, _ := db.SystemUser.Where(db.SystemUser.Account.Eq(req.Account)).First()
	if user != nil {
		return errors.New("该账号已被使用")
	}

	salt := utils.RandStr(8, utils.AllDict)
	if err := db.SystemUser.Create(&model.SystemUser{
		Account:  req.Account,
		Password: utils.MD5AndSalt(req.Password, salt),
		Salt:     salt,
		Admin:    req.Admin,
		Amount:   req.Amount,
	}); err != nil {
		log.Error("新增用户失败", zap.Error(err))
		return errors.New("操作失败")
	}
	return nil
}
