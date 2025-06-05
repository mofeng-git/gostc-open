package service

import (
	"server/pkg/utils"
	"server/repository"
	"server/service/common/cache"
	"sort"
)

func (service *service) NodeObs() (result []NodeObsItem) {
	db, _, _ := repository.Get("")
	nodes, _ := db.GostNode.Select(
		db.GostNode.Code,
		db.GostNode.Name,
	).Find()
	var nodeObsMap = make(map[string]NodeObsItem)
	for _, dateOnly := range cache.MONTH_DATEONLY_LIST {
		for _, node := range nodes {
			obsInfo := cache.GetNodeObs(dateOnly, node.Code)
			obs := nodeObsMap[node.Code]
			obs.Code = node.Code
			obs.Name = node.Name
			obs.Online = utils.TrinaryOperation(cache.GetNodeOnline(node.Code), 1, 2)
			obs.InputBytes += obsInfo.InputBytes
			obs.OutputBytes += obsInfo.OutputBytes
			nodeObsMap[node.Code] = obs
		}
	}

	var validNodeObsList nodeObsSortable
	for _, obs := range nodeObsMap {
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
