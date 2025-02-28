package service

type NodeListReq struct {
	Code       string `binding:"required" json:"code"`
	TunnelType int    `binding:"required" json:"tunnelType"`
}

type NodeListItem struct {
	Code    string   `json:"code"`
	Name    string   `json:"name"`
	Remark  string   `json:"remark"`
	Web     int      `json:"web"`
	Tunnel  int      `json:"tunnel"`
	Forward int      `json:"forward"`
	Rules   []string `json:"rules"`
	Tags    []string `json:"tags"`
}

func (service *service) NodeList(req NodeListReq) (list []NodeListItem, err error) {
	//var cfg model.GostClientConfig
	//if db.Where("code = ?", req.Code).First(&cfg).RowsAffected == 0 {
	//	return nil, errors.New("套餐配置错误")
	//}
	//var nodes []model.GostNode
	//var where = db.Where("code in ?", cfg.GetNodes())
	//switch req.TunnelType {
	//case 1:
	//	where = where.Where("web = 1")
	//case 2:
	//	where = where.Where("forward = 1")
	//case 3:
	//	where = where.Where("tunnel = 1")
	//}
	//db.Where(where).Find(&nodes)
	//if len(nodes) == 0 {
	//	return nil, errors.New("该套餐无有效节点")
	//}
	//
	//for _, node := range nodes {
	//	var ruleNames []string
	//	for _, rule := range node.GetRules() {
	//		ruleNames = append(ruleNames, node_rule.RuleMap[rule].Name())
	//	}
	//	list = append(list, NodeListItem{
	//		Code:    node.Code,
	//		Name:    node.Name,
	//		Remark:  node.Remark,
	//		Web:     node.Web,
	//		Tunnel:  node.Tunnel,
	//		Forward: node.Forward,
	//		Rules:   ruleNames,
	//		Tags:    node.GetTags(),
	//	})
	//}
	return list, nil
}
