package service

import (
	"errors"
	"go.uber.org/zap"
	"server/model"
	"server/repository"
)

type DeleteReq struct {
	Code string `binding:"required" json:"code" label:"编号"`
}

func (service *service) Delete(req DeleteReq) error {
	db, _, log := repository.Get("")
	var cfg model.GostNodeConfig
	if db.Where("code = ?", req.Code).First(&cfg).RowsAffected == 0 {
		return nil
	}

	if err := db.Delete(&cfg).Error; err != nil {
		log.Error("删除套餐配置失败", zap.Error(err))
		return errors.New("操作失败")
	}
	return nil
}
