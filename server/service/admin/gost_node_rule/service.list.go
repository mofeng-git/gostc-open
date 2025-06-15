package service

import "server/service/common/node_rule"

type Item struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (*service) List() (list []Item) {
	for _, v := range node_rule.Registry.GetRules() {
		list = append(list, Item{
			Code:        v.Code(),
			Name:        v.Name(),
			Description: v.Description(),
		})
	}
	return list
}
