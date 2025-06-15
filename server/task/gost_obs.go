package task

import (
	"server/model"
	"server/repository"
	"server/repository/cache"
	"time"
)

func gostObs() {
	db, _, _ := repository.Get("")
	now := time.Now()
	prevDate := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local).Format(time.DateOnly)
	var obsList []*model.GostObs
	{
		var codes []string
		_ = db.GostClientHost.Pluck(db.GostClientHost.Code, &codes)
		for _, code := range codes {
			info := cache.GetTunnelObs(prevDate, code)
			obsList = append(obsList, &model.GostObs{
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
		_ = db.GostClientForward.Pluck(db.GostClientForward.Code, &codes)
		for _, code := range codes {
			info := cache.GetTunnelObs(prevDate, code)
			obsList = append(obsList, &model.GostObs{
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
		_ = db.GostClientTunnel.Pluck(db.GostClientTunnel.Code, &codes)
		for _, code := range codes {
			info := cache.GetTunnelObs(prevDate, code)
			obsList = append(obsList, &model.GostObs{
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
		_ = db.GostClient.Pluck(db.GostClient.Code, &codes)
		for _, code := range codes {
			info := cache.GetClientObs(prevDate, code)
			obsList = append(obsList, &model.GostObs{
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
		_ = db.GostNode.Pluck(db.GostNode.Code, &codes)
		for _, code := range codes {
			info := cache.GetNodeObs(prevDate, code)
			obsList = append(obsList, &model.GostObs{
				OriginKind:  model.GOST_OBS_ORIGIN_KIND_NODE,
				OriginCode:  code,
				EventDate:   prevDate,
				InputBytes:  info.InputBytes,
				OutputBytes: info.OutputBytes,
			})
		}
	}
	{
		var codes []string
		_ = db.SystemUser.Pluck(db.SystemUser.Code, &codes)
		for _, code := range codes {
			info := cache.GetUserObs(prevDate, code)
			obsList = append(obsList, &model.GostObs{
				OriginKind:  model.GOST_OBS_ORIGIN_KIND_USER,
				OriginCode:  code,
				EventDate:   prevDate,
				InputBytes:  info.InputBytes,
				OutputBytes: info.OutputBytes,
			})
		}
	}
	_ = db.GostObs.CreateInBatches(obsList, 1000)

	{
		nodes, _ := db.GostNode.Select(
			db.GostNode.Code,
			db.GostNode.LimitResetIndex,
		).Where().Find()
		for _, node := range nodes {
			cache.RefreshNodeObsLimit(node.Code, node.LimitResetIndex)
		}
	}
}
