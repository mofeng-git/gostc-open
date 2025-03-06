package service

import (
	"errors"
	"go.uber.org/zap"
	"server/repository"
)

type DeleteReq struct {
	Code string `binding:"required" json:"code" label:"编号"`
}

func (service *service) Delete(req DeleteReq) error {
	db, _, log := repository.Get("")
	cfg, _ := db.GostNodeConfig.Where(db.GostNodeConfig.Code.Eq(req.Code)).First()
	if cfg == nil {
		return nil
	}

	if _, err := db.GostNodeConfig.Where(db.GostNodeConfig.Code.Eq(cfg.Code)).Delete(); err != nil {
		log.Error("删除套餐配置失败", zap.Error(err))
		return errors.New("操作失败")
	}
	return nil
}
