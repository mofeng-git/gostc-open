package service

import (
	"errors"
	"server/repository"
)

type DeleteReq struct {
	Code string `binding:"required" json:"code" label:"编号"`
}

func (service *service) Delete(req DeleteReq) error {
	db, _, _ := repository.Get("")
	if _, err := db.SystemNotice.Where(db.SystemNotice.Code.Eq(req.Code)).Delete(); err != nil {
		return errors.New("操作失败")
	}
	return nil
}
