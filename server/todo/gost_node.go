package todo

import (
	cache2 "github.com/patrickmn/go-cache"
	"server/repository"
	"server/repository/cache"
)

func gostNode() {
	db, _, _ := repository.Get("")
	nodes, _ := db.GostNode.Select(
		db.GostNode.Code,
		db.GostNode.LimitTotal,
		db.GostNode.LimitKind,
		db.GostNode.LimitResetIndex,
	).Find()
	for _, node := range nodes {
		cache.SetNodeInfo(cache.NodeInfo{
			Code:            node.Code,
			LimitResetIndex: node.LimitResetIndex,
			LimitTotal:      node.LimitTotal,
			LimitKind:       node.LimitKind,
		})
		cache.SetNodeOnline(node.Code, false, cache2.NoExpiration)
	}
}
