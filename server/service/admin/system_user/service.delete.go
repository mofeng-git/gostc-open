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
	user, _ := db.SystemUser.Where(db.SystemUser.Code.Eq(req.Code)).First()
	if user == nil {
		return nil
	}

	if _, err := db.SystemUser.Where(db.SystemUser.Code.Eq(user.Code)).Delete(); err != nil {
		log.Error("删除用户失败", zap.Error(err))
		return errors.New("操作失败")
	}
	return nil
}
