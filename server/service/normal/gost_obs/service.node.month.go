package service

import (
	"server/pkg/utils"
	"server/service/common/cache"
	"time"
)

type NodeMonthReq struct {
	Code  string `binding:"required" json:"code"`
	Start string `json:"start"`
	End   string `json:"end"`
}

func (service *service) NodeMonth(req NodeMonthReq) (result []Item) {
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
		summary := cache.GetNodeObs(date, req.Code)
		result = append(result, Item{
			Date: date,
			In:   summary.InputBytes,
			Out:  summary.OutputBytes,
		})
	}
	return result
}
