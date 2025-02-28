package service

import (
	"server/model"
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
	var nodes []model.GostNode
	var where = db
	switch req.Bind {
	case 1:
		// 用户节点
		where = where.Where("code IN (?)", db.Model(&model.GostNodeBind{}).Select("node_code"))
	case 2:
		// 系统节点
		where = where.Where("code NOT IN (?)", db.Model(&model.GostNodeBind{}).Select("node_code"))
	}
	db.Where(where).Select("code, name").Order("index_value asc").Order("id desc").Find(&nodes)
	for _, node := range nodes {
		list = append(list, ListItem{
			Code:   node.Code,
			Name:   node.Name,
			Online: utils.TrinaryOperation(cache.GetNodeOnline(node.Code), 1, 2),
		})
	}
	return list
}
