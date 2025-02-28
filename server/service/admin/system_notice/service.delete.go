package service

import (
	"errors"
	"server/model"
	"server/repository"
)

type DeleteReq struct {
	Code string `binding:"required" json:"code" label:"编号"`
}

func (service *service) Delete(req DeleteReq) error {
	db, _, _ := repository.Get("")
	var notice model.SystemNotice
	if db.Where("code = ?", req.Code).First(&notice).RowsAffected == 0 {
		return nil
	}

	if err := db.Delete(&notice).Error; err != nil {
		return errors.New("操作失败")
	}
	return nil
}
