package service

import (
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/repository/cache"
)

type ListItem struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Online int    `json:"online"`
}

func (service *service) List(claims jwt.Claims) (list []ListItem) {
	db, _, _ := repository.Get("")
	clients, _ := db.GostClient.Where(db.GostClient.UserCode.Eq(claims.Code)).Order(db.GostClient.Id.Desc()).Find()
	for _, client := range clients {
		list = append(list, ListItem{
			Code:   client.Code,
			Name:   client.Name,
			Online: utils.TrinaryOperation(cache.GetClientOnline(client.Code), 1, 2),
		})
	}
	return list
}
