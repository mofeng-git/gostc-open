package service

import (
	"server/repository"
	"server/repository/cache"
	"sort"
)

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
