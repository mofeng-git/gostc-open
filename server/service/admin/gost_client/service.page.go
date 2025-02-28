package service

import (
	"server/model"
	"server/pkg/bean"
	"server/pkg/utils"
	"server/repository"
	"server/service/common/cache"
	"time"
)

type PageReq struct {
	bean.PageParam
	Account string `json:"account"`
}

type Item struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Ip          string `json:"ip"`
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
	var clients []model.GostClient
	var where = db
	if req.Account != "" {
		where = where.Where("user_code in (?)",
			db.Model(&model.SystemUser{}).Where("account like ?", "%"+req.Account+"%").Select("code"),
		)
	}
	db.Where(where).Model(&clients).Count(&total)
	db.Preload("User").Where(where).Order("id desc").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&clients)
	for _, client := range clients {
		obsInfo := cache.GetClientObsDateRange(cache.MONTH_DATEONLY_LIST, client.Code)
		list = append(list, Item{
			Code:        client.Code,
			Name:        client.Name,
			Key:         client.Key,
			Ip:          cache.GetClientIp(client.Code),
			Online:      utils.TrinaryOperation(cache.GetClientOnline(client.Code), 1, 2),
			LastTime:    cache.GetClientLastTime(client.Code),
			Version:     cache.GetClientVersion(client.Code),
			CreatedAt:   client.CreatedAt.Format(time.DateTime),
			InputBytes:  obsInfo.InputBytes,
			OutputBytes: obsInfo.OutputBytes,
			UserCode:    client.UserCode,
			UserAccount: client.User.Account,
		})
	}
	return list, total
}
