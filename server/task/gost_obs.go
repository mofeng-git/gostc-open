package task

import (
	"server/bootstrap"
	"server/model"
	"server/repository"
	"server/service/common/cache"
	"time"
)

func init() {
	bootstrap.AddCron("0 0 * * *", func() {
		db, _, _ := repository.Get("")
		now := time.Now()
		prevDate := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local).Format(time.DateOnly)
		var obsList []model.GostObs
		{
			var codes []string
			db.Model(&model.GostClientHost{}).Pluck("code", &codes)
			for _, code := range codes {
				info := cache.GetTunnelObs(prevDate, code)
				obsList = append(obsList, model.GostObs{
					OriginKind:  model.GOST_OBS_ORIGIN_KIND_TUNNEL,
					OriginCode:  code,
					EventDate:   prevDate,
					InputBytes:  info.InputBytes,
					OutputBytes: info.OutputBytes,
				})
			}
		}
		{
			var codes []string
			db.Model(&model.GostClientForward{}).Pluck("code", &codes)
			for _, code := range codes {
				info := cache.GetTunnelObs(prevDate, code)
				obsList = append(obsList, model.GostObs{
					OriginKind:  model.GOST_OBS_ORIGIN_KIND_TUNNEL,
					OriginCode:  code,
					EventDate:   prevDate,
					InputBytes:  info.InputBytes,
					OutputBytes: info.OutputBytes,
				})
			}
		}
		{
			var codes []string
			db.Model(&model.GostClientTunnel{}).Pluck("code", &codes)
			for _, code := range codes {
				info := cache.GetTunnelObs(prevDate, code)
				obsList = append(obsList, model.GostObs{
					OriginKind:  model.GOST_OBS_ORIGIN_KIND_TUNNEL,
					OriginCode:  code,
					EventDate:   prevDate,
					InputBytes:  info.InputBytes,
					OutputBytes: info.OutputBytes,
				})
			}
		}
		{
			var codes []string
			db.Model(&model.GostClient{}).Pluck("code", &codes)
			for _, code := range codes {
				info := cache.GetClientObs(prevDate, code)
				obsList = append(obsList, model.GostObs{
					OriginKind:  model.GOST_OBS_ORIGIN_KIND_CLIENT,
					OriginCode:  code,
					EventDate:   prevDate,
					InputBytes:  info.InputBytes,
					OutputBytes: info.OutputBytes,
				})
			}
		}
		{
			var codes []string
			db.Model(&model.GostNode{}).Pluck("code", &codes)
			for _, code := range codes {
				info := cache.GetNodeObs(prevDate, code)
				obsList = append(obsList, model.GostObs{
					OriginKind:  model.GOST_OBS_ORIGIN_KIND_NODE,
					OriginCode:  code,
					EventDate:   prevDate,
					InputBytes:  info.InputBytes,
					OutputBytes: info.OutputBytes,
				})
			}
		}
		db.CreateInBatches(&obsList, 1000)
	})
}
