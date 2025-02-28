package service

import (
	"errors"
	"server/model"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"time"
)

type UserInfoResp struct {
	Account       string `json:"account"`
	Amount        string `json:"amount"`
	CheckinAmount string `json:"checkinAmount"`
	Admin         int    `json:"admin"`
	Otp           int    `json:"otp"`
	CreatedAt     string `json:"createdAt"`
}

func (service *service) UserInfo(claims jwt.Claims) (result UserInfoResp, err error) {
	db, _, _ := repository.Get("")
	var user model.SystemUser
	if db.Preload("BindQQ").Where("code = ?", claims.Code).First(&user).RowsAffected == 0 {
		return result, errors.New("未查询到账户信息")
	}
	var checkIn model.SystemUserCheckin
	db.Where("user_code = ? AND event_date = ?", user.Code, time.Now().Format(time.DateOnly)).First(&checkIn)
	return UserInfoResp{
		Account:       user.Account,
		Amount:        user.Amount.String(),
		Admin:         user.Admin,
		Otp:           utils.TrinaryOperation(user.OtpKey == "", 2, 1),
		CheckinAmount: checkIn.Amount.String(),
		CreatedAt:     user.CreatedAt.Format(time.DateTime),
	}, nil
}
