package service

import (
	"errors"
	"go.uber.org/zap"
	"server/repository"
	"server/repository/query"
)

type DeleteReq struct {
	Code string `binding:"required" json:"code" label:"编号"`
}

func (service *service) Delete(req DeleteReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(req.Code)).First()
		if user == nil {
			return nil
		}

		client, _ := tx.GostClient.Where(tx.GostClient.UserCode.Eq(req.Code)).First()
		if client != nil {
			return errors.New("请先删除该用户所有的客户端")
		}

		if _, err := tx.SystemUser.Where(tx.SystemUser.Code.Eq(user.Code)).Delete(); err != nil {
			log.Error("删除用户失败", zap.Error(err))
			return errors.New("操作失败")
		}
		return nil
	})
}
