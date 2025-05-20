package todo

import (
	"server/model"
	"server/repository"
	"server/service/common/cache"
)

// 从数据库回写流量数据到cache
func gostObsWriteBack() {
	db, _, _ := repository.Get("")
	obsList, _ := db.GostObs.Where(db.GostObs.EventDate.In(cache.MONTH_DATEONLY_LIST[:29]...)).Find()
	for _, obs := range obsList {
		if obs.InputBytes == 0 && obs.OutputBytes == 0 {
			continue
		}
		switch obs.OriginKind {
		case model.GOST_OBS_ORIGIN_KIND_CLIENT:
			cache.OverflowClientObs(obs.EventDate, obs.OriginCode, cache.TunnelObs{
				InputBytes:  obs.InputBytes,
				OutputBytes: obs.OutputBytes,
			})
		case model.GOST_OBS_ORIGIN_KIND_TUNNEL:
			cache.OverflowTunnelObs(obs.EventDate, obs.OriginCode, cache.TunnelObs{
				InputBytes:  obs.InputBytes,
				OutputBytes: obs.OutputBytes,
			})
		case model.GOST_OBS_ORIGIN_KIND_NODE:
			cache.OverflowNodeObs(obs.EventDate, obs.OriginCode, cache.TunnelObs{
				InputBytes:  obs.InputBytes,
				OutputBytes: obs.OutputBytes,
			})
		case model.GOST_OBS_ORIGIN_KIND_USER:
			cache.OverflowUserObs(obs.EventDate, obs.OriginCode, cache.TunnelObs{
				InputBytes:  obs.InputBytes,
				OutputBytes: obs.OutputBytes,
			})
		}
	}
}
