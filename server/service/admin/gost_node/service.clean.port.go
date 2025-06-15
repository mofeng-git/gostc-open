package service

import (
	"server/repository"
	"server/repository/cache"
	"server/service/common/node_port"
)

type CleanPortReq struct {
	Code string `binding:"required" json:"code"`
}

func (service *service) CleanPort(req CleanPortReq) error {
	db, _, _ := repository.Get("")
	node, _ := db.GostNode.Where(db.GostNode.Code.Eq(req.Code)).First()
	if node == nil {
		return nil
	}
	cache.DelNodePortUse(node.Code)
	node_port.Arrange(db, node.Code)
	return nil
}
