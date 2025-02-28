package service

import (
	"server/model"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/service/common/cache"
)

type ListItem struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Online int    `json:"online"`
}

func (service *service) List(claims jwt.Claims) (list []ListItem) {
	db, _, _ := repository.Get("")
	var clients []model.GostClient
	var where = db.Where("user_code = ?", claims.Code)
	db.Where(where).Order("id desc").Find(&clients)
	for _, client := range clients {
		list = append(list, ListItem{
			Code:   client.Code,
			Name:   client.Name,
			Online: utils.TrinaryOperation(cache.GetClientOnline(client.Code), 1, 2),
		})
	}
	return list
}
