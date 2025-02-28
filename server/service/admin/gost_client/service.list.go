package service

import (
	"server/model"
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
	var clients []model.GostClient
	var where = db
	if req.UserCode != "" {
		where = where.Where("user_code = ?", req.UserCode)
	}
	db.Preload("User").Where(where).Order("id desc").Find(&clients)
	for _, client := range clients {
		list = append(list, ListItem{
			Code: client.Code,
			Name: client.Name,
		})
	}
	return list
}
