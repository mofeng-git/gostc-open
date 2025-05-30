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
	Email         string `json:"email"`
}

func (service *service) UserInfo(claims jwt.Claims) (result UserInfoResp, err error) {
	db, _, _ := repository.Get("")
	user, _ := db.SystemUser.Preload(db.SystemUser.BindEmail).Where(db.SystemUser.Code.Eq(claims.Code)).First()
	if user == nil {
		return result, errors.New("未查询到账户信息")
	}
	checkIn, _ := db.SystemUserCheckin.Where(
		db.SystemUserCheckin.UserCode.Eq(user.Code),
		db.SystemUserCheckin.EventDate.Eq(time.Now().Format(time.DateOnly)),
	).First()
	if checkIn == nil {
		checkIn = &model.SystemUserCheckin{}
	}

	return UserInfoResp{
		Account:       user.Account,
		Amount:        user.Amount.String(),
		Admin:         user.Admin,
		Otp:           utils.TrinaryOperation(user.OtpKey == "", 2, 1),
		CheckinAmount: checkIn.Amount.String(),
		CreatedAt:     user.CreatedAt.Format(time.DateTime),
		Email:         user.BindEmail.Email,
	}, nil
}
