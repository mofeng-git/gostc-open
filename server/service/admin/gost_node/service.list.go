package service

import (
	"gorm.io/gen"
	"server/pkg/utils"
	"server/repository"
	"server/service/common/cache"
)

type ListReq struct {
	Bind int `json:"bind"`
}

type ListItem struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Online int    `json:"online"`
}

func (service *service) List(req ListReq) (list []ListItem) {
	db, _, _ := repository.Get("")
	var where []gen.Condition

	var nodeCodes []string
	_ = db.GostNodeBind.Pluck(db.GostNodeBind.NodeCode, &nodeCodes)
	switch req.Bind {
	case 1:
		// 用户节点
		where = append(where, db.GostNode.Code.In(nodeCodes...))
	case 2:
		// 系统节点
		where = append(where, db.GostNode.Code.NotIn(nodeCodes...))
	}

	nodes, _ := db.GostNode.Where(where...).Select(db.GostNode.Code, db.GostNode.Name).Order(db.GostNode.IndexValue.Asc(), db.GostNode.Id.Asc()).Find()
	for _, node := range nodes {
		list = append(list, ListItem{
			Code:   node.Code,
			Name:   node.Name,
			Online: utils.TrinaryOperation(cache.GetNodeOnline(node.Code), 1, 2),
		})
	}
	return list
}
