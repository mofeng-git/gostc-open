package todo

import (
	"github.com/shopspring/decimal"
	"server/bootstrap"
	"server/model"
	"server/pkg/utils"
	"server/repository"
)

func init() {
	bootstrap.AddTodo(func() {
		db, _, _ := repository.Get("")
		salt := utils.RandStr(8, utils.AllDict)
		db.Create(&model.SystemUser{
			Base: model.Base{
				Id:   1,
				Code: "1",
			},
			Account:  "admin",
			Password: utils.MD5AndSalt("admin", salt),
			Salt:     salt,
			OtpKey:   "",
			Admin:    1,
			Amount:   decimal.NewFromInt(1000),
		})
	})
}
