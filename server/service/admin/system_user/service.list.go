package service

import (
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
	users, _ := db.SystemUser.Where(db.SystemUser.Account.Like("%"+req.Account+"%")).Select(
		db.SystemUser.Code,
		db.SystemUser.Account,
	).Order(db.SystemUser.Id.Desc()).Find()
	for _, user := range users {
		list = append(list, ListItem{
			Code:    user.Code,
			Account: user.Account,
		})
	}
	return list
}
