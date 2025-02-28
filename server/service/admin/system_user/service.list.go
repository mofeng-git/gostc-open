package service

import (
	"server/model"
	"server/repository"
)

type ListReq struct {
	Account string `binding:"required" json:"account"`
}

type ListItem struct {
	Code    string `json:"code"`
	Account string `json:"account"`
}

func (service *service) List(req ListReq) (list []ListItem) {
	db, _, _ := repository.Get("")
	var users []model.SystemUser
	db.Select("code, account").Where("account like ?", "%"+req.Account+"%").Order("id desc").Find(&users)
	for _, user := range users {
		list = append(list, ListItem{
			Code:    user.Code,
			Account: user.Account,
		})
	}
	return list
}
