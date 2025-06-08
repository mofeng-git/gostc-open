package service

type CheckInReq struct {
	QQ string `binding:"required" json:"QQ"`
}

func (service *service) Checkin(req CheckInReq) (score string, err error) {
	//var cfg model.SystemConfigBase
	//cache.GetSystemConfigBase(&cfg)
	//if cfg.CheckIn != "1" {
	//	return "0", errors.New("未启用签到功能")
	//}
	//
	//db, _, log := repository.Get("")
	//if err := db.Transaction(func(tx *query.Query) error {
	//	bindQQ, _ := tx.SystemUserQQ.Where(tx.SystemUserQQ.QQ.Eq(req.QQ)).First()
	//	if bindQQ == nil || bindQQ.QQ == "" {
	//		return errors.New("未绑定账号，请在网站通过绑定码，绑定账号")
	//	}
	//
	//	user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(bindQQ.UserCode)).First()
	//	if user == nil {
	//		return errors.New("账号错误")
	//	}
	//
	//	checkin, _ := tx.SystemUserCheckin.Where(
	//		tx.SystemUserCheckin.UserCode.Eq(user.Code),
	//		tx.SystemUserCheckin.EventDate.Eq(time.Now().Format(time.DateOnly)),
	//	).First()
	//	if checkin != nil {
	//		return errors.New("已签到")
	//	}
	//
	//	// 0-4 +6
	//	amount := decimal.NewFromInt(int64(utils.RandNum(5) + 6))
	//	user.Amount = user.Amount.Add(amount)
	//	if err := tx.SystemUser.Save(user); err != nil {
	//		log.Error("签到失败", zap.Error(err))
	//		return errors.New("签到失败")
	//	}
	//	if err := tx.SystemUserCheckin.Create(&model.SystemUserCheckin{
	//		UserCode:  user.Code,
	//		Account:   user.Account,
	//		EventDate: time.Now().Format(time.DateOnly),
	//		Amount:    amount,
	//	}); err != nil {
	//		log.Error("签到失败", zap.Error(err))
	//		return errors.New("签到失败")
	//	}
	//	score = amount.String()
	//	return nil
	//}); err != nil {
	//	return "0", err
	//}
	return score, nil
}
