package service

import (
	"gorm.io/gen"
	"server/pkg/bean"
	"server/pkg/utils"
	"server/repository"
	cache2 "server/repository/cache"
	"time"
)

type PageReq struct {
	bean.PageParam
	Name    string `json:"name"`
	Account string `json:"account"`
}

type Item struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Online      int    `json:"online"`
	LastTime    string `json:"lastTime"`
	Version     string `json:"version"`
	CreatedAt   string `json:"createdAt"`
	InputBytes  int64  `json:"inputBytes"`
	OutputBytes int64  `json:"outputBytes"`

	UserCode    string `json:"userCode"`
	UserAccount string `json:"userAccount"`
}

func (service *service) Page(req PageReq) (list []Item, total int64) {
	db, _, _ := repository.Get("")

	var where []gen.Condition
	if req.Account != "" {
		var userCodes []string
		_ = db.SystemUser.Where(db.SystemUser.Account.Like("%"+req.Account+"%")).Pluck(db.SystemUser.Code, &userCodes)
		where = append(where, db.GostClient.UserCode.In(userCodes...))
	}

	if req.Name != "" {
		where = append(where, db.GostClient.Name.Like("%"+req.Name+"%"))
	}

	clients, total, _ := db.GostClient.
		Preload(db.GostClient.User).
		Where(where...).
		Order(db.GostClient.Id.Desc()).
		FindByPage(req.GetOffset(), req.GetLimit())
	for _, client := range clients {
		obsInfo := cache2.GetClientObsDateRange(cache2.MONTH_DATEONLY_LIST, client.Code)
		list = append(list, Item{
			Code:        client.Code,
			Name:        client.Name,
			Key:         client.Key,
			Online:      utils.TrinaryOperation(cache2.GetClientOnline(client.Code), 1, 2),
			LastTime:    cache2.GetClientLastTime(client.Code),
			Version:     cache2.GetClientVersion(client.Code),
			CreatedAt:   client.CreatedAt.Format(time.DateTime),
			InputBytes:  obsInfo.InputBytes,
			OutputBytes: obsInfo.OutputBytes,
			UserCode:    client.UserCode,
			UserAccount: client.User.Account,
		})
	}
	return list, total
}
