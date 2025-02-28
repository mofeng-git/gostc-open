package service

type UpdateReq struct {
	Code         string   `binding:"required" json:"code" label:"编号"`
	Name         string   `binding:"required" json:"name"`
	ChargingType int      `binding:"required" json:"chargingType" label:"计费类型"`
	Cycle        int      `json:"cycle" label:"计费周期(天)"`
	Amount       string   `json:"amount" label:"计费金额"`
	Limiter      int      `json:"limiter" label:"速率"`
	RLimiter     int      `json:"rLimiter" label:"并发数"`
	CLimiter     int      `json:"cLimiter" label:"连接数"`
	OnlyChina    int      `json:"onlyChina" label:"仅中国大陆可用"`
	Nodes        []string `json:"nodes" label:"可使用的节点"`
	ExpAt        int64    `json:"expAt"`
}

func (service *service) Update(req UpdateReq) (err error) {
	//var amount decimal.Decimal
	//switch req.ChargingType {
	//case model.GOST_CONFIG_CHARGING_ONLY_ONCE, model.GOST_CONFIG_CHARGING_CUCLE_DAY:
	//	amount, err = decimal.NewFromString(req.Amount)
	//	if err != nil {
	//		return fmt.Errorf("金额错误，%v", err)
	//	}
	//default:
	//	req.ChargingType = model.GOST_CONFIG_CHARGING_FREE
	//}
	//
	//var cfg model.GostClientConfig
	//if db.Where("code = ?", req.Code).First(&cfg).RowsAffected == 0 {
	//	return errors.New("数据不存在")
	//}
	//
	//cfg.Name = req.Name
	//cfg.ChargingType = req.ChargingType
	//cfg.Cycle = req.Cycle
	//cfg.Amount = amount
	//cfg.Limiter = req.Limiter
	//cfg.RLimiter = req.RLimiter
	//cfg.CLimiter = req.CLimiter
	//cfg.OnlyChina = req.OnlyChina
	//cfg.Nodes = strings.Join(req.Nodes, ",")
	//cfg.ExpAt = req.ExpAt
	//if err := db.Save(&cfg).Error; err != nil {
	//	log.Error("修改用户套餐配置失败", zap.Error(err))
	//	return errors.New("操作失败")
	//}
	//if cfg.TunnelCode != "" {
	//	switch cfg.TunnelType {
	//	case model.GOST_TUNNEL_TYPE_HOST:
	//		gost_engine.ClientHostConfig(db, cfg.TunnelCode)
	//	case model.GOST_TUNNEL_TYPE_FORWARD:
	//		gost_engine.ClientForwardConfig(db, cfg.TunnelCode)
	//	case model.GOST_TUNNEL_TYPE_TUNNEL:
	//		gost_engine.ClientTunnelConfig(db, cfg.TunnelCode)
	//	}
	//}
	return nil
}
