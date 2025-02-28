package service

import (
	"github.com/shopspring/decimal"
	"server/model"
	"server/pkg/bean"
	"server/repository"
	"server/service/common/cache"
	"time"
)

type PageReq struct {
	bean.PageParam
	Account string `json:"account"`
	Admin   int    `json:"admin"`
}

type Item struct {
	Code        string          `json:"code"`
	Account     string          `json:"account"`
	Admin       int             `json:"admin"`
	Amount      decimal.Decimal `json:"amount"`
	CreatedAt   string          `json:"createdAt"`
	InputBytes  int64           `json:"inputBytes"`
	OutputBytes int64           `json:"outputBytes"`
}

func (service *service) Page(req PageReq) (list []Item, total int64) {
	db, _, _ := repository.Get("")
	var users []model.SystemUser
	var where = db
	if req.Account != "" {
		where = where.Where("account = ?", req.Account)
	}
	if req.Admin > 0 {
		where = where.Where("admin = ?", req.Admin)
	}

	db.Where(where).Model(&users).Count(&total)
	db.Where(where).Preload("BindQQ").Order("id desc").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&users)
	for _, user := range users {
		obsInfo := cache.GetUserObsDateRange(cache.MONTH_DATEONLY_LIST, user.Code)
		list = append(list, Item{
			Code:        user.Code,
			Account:     user.Account,
			Admin:       user.Admin,
			Amount:      user.Amount,
			CreatedAt:   user.CreatedAt.Format(time.DateTime),
			InputBytes:  obsInfo.InputBytes,
			OutputBytes: obsInfo.OutputBytes,
		})
	}
	return list, total
}
