package service

import (
	"server/pkg/jwt"
	"server/repository"
	"server/repository/cache"
	"sort"
	"time"
)

func (service *service) ClientHostObsDate(claims jwt.Claims, date string) (result []TunnelObsItem) {
	dateOnly, err := time.ParseInLocation(time.DateOnly, date, time.Local)
	if err != nil {
		dateOnly = time.Now()
	}
	var dateOnlyString = dateOnly.Format(time.DateOnly)

	db, _, _ := repository.Get("")
	hosts, _ := db.GostClientHost.Select(
		db.GostClientHost.Code,
		db.GostClientHost.Name,
	).Where(db.GostClientHost.UserCode.Eq(claims.Code)).Find()
	var tunnelObsMap = make(map[string]TunnelObsItem)
	for _, item := range hosts {
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
