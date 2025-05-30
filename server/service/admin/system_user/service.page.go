package service

import (
	"github.com/shopspring/decimal"
	"gorm.io/gen"
	"server/pkg/bean"
	"server/repository"
	"server/service/common/cache"
	"time"
)

type PageReq struct {
	bean.PageParam
	Account string `json:"account"`
	Admin   int    `json:"admin"`
	Email   string `json:"email"`
}

type Item struct {
	Code        string          `json:"code"`
	Account     string          `json:"account"`
	Admin       int             `json:"admin"`
	Amount      decimal.Decimal `json:"amount"`
	Email       string          `json:"email"`
	CreatedAt   string          `json:"createdAt"`
	InputBytes  int64           `json:"inputBytes"`
	OutputBytes int64           `json:"outputBytes"`
}

func (service *service) Page(req PageReq) (list []Item, total int64) {
	db, _, _ := repository.Get("")
	var where []gen.Condition
	if req.Account != "" {
		where = append(where, db.SystemUser.Account.Like("%"+req.Account+"%"))
	}
	if req.Admin > 0 {
		where = append(where, db.SystemUser.Admin.Eq(req.Admin))
	}
	users, total, _ := db.SystemUser.Preload(db.SystemUser.BindEmail).Where(where...).Order(db.SystemUser.Id.Desc()).FindByPage(req.GetOffset(), req.GetLimit())
	for _, user := range users {
		obsInfo := cache.GetUserObsDateRange(cache.MONTH_DATEONLY_LIST, user.Code)
		list = append(list, Item{
			Code:        user.Code,
			Account:     user.Account,
			Admin:       user.Admin,
			Amount:      user.Amount,
			Email:       user.BindEmail.Email,
			CreatedAt:   user.CreatedAt.Format(time.DateTime),
			InputBytes:  obsInfo.InputBytes,
			OutputBytes: obsInfo.OutputBytes,
		})
	}
	return list, total
}
