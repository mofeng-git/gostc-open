package service

import (
	"server/repository"
	"server/service/common/cache"
	"sort"
	"time"
)

func (service *service) ClientForwardObsDate(date string) (result []TunnelObsItem) {
	dateOnly, err := time.ParseInLocation(time.DateOnly, date, time.Local)
	if err != nil {
		dateOnly = time.Now()
	}
	var dateOnlyString = dateOnly.Format(time.DateOnly)

	db, _, _ := repository.Get("")
	forwards, _ := db.GostClientForward.Select(
		db.GostClientForward.Code,
		db.GostClientForward.Name,
	).Find()
	var tunnelObsMap = make(map[string]TunnelObsItem)
	for _, item := range forwards {
		obsInfo := cache.GetTunnelObs(dateOnlyString, item.Code)
		obs := tunnelObsMap[item.Code]
		obs.Code = item.Code
		obs.Name = item.Name
		obs.InputBytes += obsInfo.InputBytes
		obs.OutputBytes += obsInfo.OutputBytes
		tunnelObsMap[item.Code] = obs
	}

	var validNodeObsList tunnelObsSortable
	for _, obs := range tunnelObsMap {
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
