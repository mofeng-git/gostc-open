package service

import (
	"server/repository"
	"server/service/common/cache"
	"sort"
)

type UserObsItem struct {
	Account     string `json:"account"`
	InputBytes  int64  `json:"inputBytes"`
	OutputBytes int64  `json:"outputBytes"`
}

type userObsSortable []UserObsItem

func (u userObsSortable) Len() int {
	return len(u)
}

func (u userObsSortable) Less(i, j int) bool {
	return u[i].InputBytes+u[i].OutputBytes > u[j].InputBytes+u[j].OutputBytes
}

func (u userObsSortable) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (service *service) UserObs() (result []UserObsItem) {
	db, _, _ := repository.Get("")
	users, _ := db.SystemUser.Select(
		db.SystemUser.Code,
		db.SystemUser.Account,
	).Find()
	var userObsMap = make(map[string]UserObsItem)
	for _, dateOnly := range cache.MONTH_DATEONLY_LIST {
		for _, user := range users {
			obsInfo := cache.GetUserObs(dateOnly, user.Code)
			obs := userObsMap[user.Code]
			obs.Account = user.Account
			obs.InputBytes += obsInfo.InputBytes
			obs.OutputBytes += obsInfo.OutputBytes
			userObsMap[user.Code] = obs
		}
	}

	var validUserObsList userObsSortable
	for _, obs := range userObsMap {
		if obs.InputBytes > 0 && obs.OutputBytes > 0 {
			validUserObsList = append(validUserObsList, obs)
		}
	}
	sort.Sort(validUserObsList)
	if len(validUserObsList) >= 30 {
		return validUserObsList[:30]
	}
	return validUserObsList
}
