package service

import (
	"errors"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"server/pkg/utils"
	"server/repository"
)

type UpdateReq struct {
	Code     string          `binding:"required" json:"code" label:"编号"`
	Account  string          `binding:"required" json:"account" label:"账号"`
	Amount   decimal.Decimal `binding:"required" json:"amount" label:"积分"`
	Password string          `json:"password" label:"密码"`
}

func (service *service) Update(req UpdateReq) error {
	db, _, log := repository.Get("")
	user, _ := db.SystemUser.Where(db.SystemUser.Code.Eq(req.Code)).First()
	if user == nil {
		return errors.New("用户不存在")
	}
	user.Account = req.Account
	user.Amount = req.Amount
	if req.Password != "" {
		user.Salt = utils.RandStr(8, utils.AllDict)
		user.Password = utils.MD5AndSalt(req.Password, user.Salt)
	}

	if err := db.SystemUser.Save(user); err != nil {
		log.Error("修改用户失败", zap.Error(err))
		return errors.New("操作失败")
	}
	return nil
}
