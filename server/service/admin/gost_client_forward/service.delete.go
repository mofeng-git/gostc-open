package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/model"
	"server/repository"
	"server/service/common/cache"
	"server/service/common/node_port"
	"server/service/gost_engine"
)

type DeleteReq struct {
	Code string `binding:"required" json:"code" label:"编号"`
}

func (service *service) Delete(req DeleteReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *gorm.DB) error {
		var forward model.GostClientForward
		if tx.Preload("Node").Where("code = ?", req.Code).First(&forward).RowsAffected == 0 {
			return nil
		}
		tx.Where("port = ? AND node_code = ?", forward.Port, forward.NodeCode).Delete(&model.GostNodePort{})
		tx.Where("tunnel_code = ?", forward.Code).Delete(&model.GostAuth{})
		if err := tx.Omit("Node").Delete(&forward).Error; err != nil {
			log.Error("删除用户端口转发失败", zap.Error(err))
			return errors.New("操作失败")
		}
		node_port.ReleasePort(forward.NodeCode, forward.Port)
		gost_engine.ClientRemoveForwardConfig(forward, forward.Node)
		cache.DelTunnelInfo(req.Code)
		return nil
	})
}
