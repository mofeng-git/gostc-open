package service

import (
	"server/pkg/jwt"
	"server/repository"
	"server/repository/cache"
	"server/service/common/node_port"
)

type CleanPortReq struct {
	Code string `binding:"required" json:"code"`
}

func (service *service) CleanPort(claims jwt.Claims, req CleanPortReq) error {
	db, _, _ := repository.Get("")
	nodeBind, _ := db.GostNodeBind.Where(db.GostNodeBind.NodeCode.Eq(req.Code), db.GostNodeBind.UserCode.Eq(claims.Code)).First()
	if nodeBind == nil {
		return nil
	}

	node, _ := db.GostNode.Where(db.GostNode.Code.Eq(req.Code)).First()
	if node == nil {
		return nil
	}
	cache.DelNodePortUse(node.Code)
	node_port.Arrange(db, node.Code)
	return nil
}
