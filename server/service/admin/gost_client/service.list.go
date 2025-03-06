package service

import (
	"gorm.io/gen"
	"server/repository"
)

type ListReq struct {
	UserCode string `json:"userCode"`
}

type ListItem struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (service *service) List(req ListReq) (list []ListItem) {
	db, _, _ := repository.Get("")

	var where []gen.Condition
	if req.UserCode != "" {
		where = append(where, db.GostClient.UserCode.Eq(req.UserCode))
	}

	clients, _ := db.GostClient.Preload(db.GostClient.User).Where(where...).Find()
	for _, client := range clients {
		list = append(list, ListItem{
			Code: client.Code,
			Name: client.Name,
		})
	}
	return list
}
