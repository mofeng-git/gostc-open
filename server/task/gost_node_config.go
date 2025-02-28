package task

func init() {
	//bootstrap.AddCron("*/10 * * * *", func() {
	//	db, _, _ := repository.Get("")
	//	_ = db.Transaction(func(tx *gorm.DB) error {
	//		where := tx.Where("charging_type = ? AND exp_at <= ? AND status = 1", model.GOST_CONFIG_CHARGING_CUCLE_DAY, time.Now().Unix())
	//		var hosts []model.GostClientHost
	//		tx.Preload("Node").Where(where).Find(&hosts)
	//		tx.Model(&model.GostClientHost{}).Where(where).Update("status", 2)
	//		for _, item := range hosts {
	//			gost_engine.ClientRemoveHostConfig(item, item.Node)
	//		}
	//		var forwards []model.GostClientForward
	//		tx.Preload("Node").Where(where).Find(&forwards)
	//		tx.Model(&model.GostClientForward{}).Where(where).Update("status", 2)
	//		for _, item := range forwards {
	//			gost_engine.ClientRemoveForwardConfig(item, item.Node)
	//		}
	//		var tunnels []model.GostClientTunnel
	//		tx.Select("code, client_code").Where(where).Find(&tunnels)
	//		tx.Model(&model.GostClientTunnel{}).Where(where).Update("status", 2)
	//		for _, item := range tunnels {
	//			gost_engine.ClientRemoveTunnelConfig(item, item.Node)
	//		}
	//		return nil
	//	})
	//})
}
