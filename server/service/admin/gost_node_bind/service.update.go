package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/model"
	"server/repository"
)

type UpdateReq struct {
	NodeCode string `binding:"required" json:"nodeCode"`
	UserCode string `json:"userCode"`
}

func (service *service) Update(req UpdateReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *gorm.DB) error {
		if req.UserCode == "" {
			tx.Where("node_code = ?", req.NodeCode).Delete(&model.GostNodeBind{})
			return nil
		}
		var bind model.GostNodeBind
		tx.Where("user_code = ? AND node_code = ?", req.UserCode, req.NodeCode).First(&bind)
		bind.UserCode = req.UserCode
		bind.NodeCode = req.NodeCode
		if err := tx.Save(&bind).Error; err != nil {
			log.Error("节点绑定用户失败", zap.Error(err))
			return errors.New("操作失败")
		}
		return nil
	})
}
