package service

import (
	"errors"
	"go.uber.org/zap"
	"server/model"
	"server/repository"
	"server/repository/query"
)

type UpdateReq struct {
	NodeCode string `binding:"required" json:"nodeCode"`
	UserCode string `json:"userCode"`
}

func (service *service) Update(req UpdateReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		if req.UserCode == "" {
			_, _ = tx.GostNodeBind.Where(tx.GostNodeBind.NodeCode.Eq(req.NodeCode)).Delete()
			return nil
		}

		bind, _ := tx.GostNodeBind.Where(
			tx.GostNodeBind.UserCode.Eq(req.UserCode),
			tx.GostNodeBind.NodeCode.Eq(req.NodeCode),
		).First()
		if bind == nil {
			bind = &model.GostNodeBind{}
		}
		bind.UserCode = req.UserCode
		bind.NodeCode = req.NodeCode
		if err := tx.GostNodeBind.Save(bind); err != nil {
			log.Error("节点绑定用户失败", zap.Error(err))
			return errors.New("操作失败")
		}
		return nil
	})
}
