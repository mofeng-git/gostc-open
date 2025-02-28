package service

import (
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/service/common/cache"
	"time"
)

type UserMonthReq struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func (service *service) UserMonth(claims jwt.Claims, req UserMonthReq) (result []Item) {
	times, ok := utils.DateFormatLayout(time.DateOnly, req.Start, req.End)
	var start, end time.Time
	if ok {
		start = times[0]
		end = times[1]
	} else {
		start = time.Now().AddDate(0, 0, -29)
		end = time.Now()
	}
	_, times2 := utils.DateRangeSplit(start, end)
	for _, date := range times2 {
		summary := cache.GetUserObs(date, claims.Code)
		result = append(result, Item{
			Date: date,
			In:   summary.InputBytes,
			Out:  summary.OutputBytes,
		})
	}
	return result
}
