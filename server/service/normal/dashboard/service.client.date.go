package service

import (
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/service/common/cache"
	"sort"
	"time"
)

func (service *service) ClientObsDate(claims jwt.Claims, date string) (result []ClientObsItem) {
	dateOnly, err := time.ParseInLocation(time.DateOnly, date, time.Local)
	if err != nil {
		dateOnly = time.Now()
	}
	var dateOnlyString = dateOnly.Format(time.DateOnly)

	db, _, _ := repository.Get("")
	forwards, _ := db.GostClient.Select(
		db.GostClient.Code,
		db.GostClient.Name,
	).Where(db.GostClient.UserCode.Eq(claims.Code)).Find()
	var clientObsMap = make(map[string]ClientObsItem)
	for _, item := range forwards {
		obsInfo := cache.GetClientObs(dateOnlyString, item.Code)
		obs := clientObsMap[item.Code]
		obs.Code = item.Code
		obs.Name = item.Name
		obs.Online = utils.TrinaryOperation(cache.GetClientOnline(item.Code), 1, 2)
		obs.InputBytes += obsInfo.InputBytes
		obs.OutputBytes += obsInfo.OutputBytes
		clientObsMap[item.Code] = obs
	}

	var validNodeObsList clientObsSortable
	for _, obs := range clientObsMap {
		if obs.InputBytes > 0 && obs.OutputBytes > 0 {
			validNodeObsList = append(validNodeObsList, obs)
		}
	}
	sort.Sort(validNodeObsList)
	if len(validNodeObsList) >= 30 {
		return validNodeObsList[:30]
	}
	return validNodeObsList
}
