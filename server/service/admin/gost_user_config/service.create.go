package service

type CreateReq struct {
	UserCode     string   `binding:"required" json:"userCode" label:"用户编号"`
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

func (service *service) Create(req CreateReq) (err error) {
	//var amount decimal.Decimal
	//switch req.ChargingType {
	//case model.GOST_CONFIG_CHARGING_ONLY_ONCE, model.GOST_CONFIG_CHARGING_CUCLE_DAY:
	//	amount, err = decimal.NewFromString(req.Amount)
	//	if err != nil {
	//		return fmt.Errorf("金额错误，%v", err)
	//	}
	//}
	//
	//if db.Where("code = ?", req.UserCode).First(&model.SystemUser{}).RowsAffected == 0 {
	//	return errors.New("用户错误")
	//}
	//
	//if err = db.Create(&model.GostClientConfig{
	//	UserCode:     req.UserCode,
	//	Name:         req.Name,
	//	ChargingType: req.ChargingType,
	//	Cycle:        req.Cycle,
	//	Amount:       amount,
	//	Limiter:      req.Limiter,
	//	RLimiter:     req.RLimiter,
	//	CLimiter:     req.CLimiter,
	//	OnlyChina:    req.OnlyChina,
	//	Nodes:        strings.Join(req.Nodes, ","),
	//	ExpAt:        req.ExpAt,
	//}).Error; err != nil {
	//	log.Error("新增用户套餐配置失败", zap.Error(err))
	//	return errors.New("操作失败")
	//}
	return nil
}
