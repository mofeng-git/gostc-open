package todo

import (
	"server/repository"
	"server/repository/cache"
)

func gostObs() {
	db, _, _ := repository.Get("")
	nodes, _ := db.GostNode.Select(
		db.GostNode.Code,
		db.GostNode.LimitResetIndex,
	).Where().Find()
	for _, node := range nodes {
		cache.RefreshNodeObsLimit(node.Code, node.LimitResetIndex)
	}
}
