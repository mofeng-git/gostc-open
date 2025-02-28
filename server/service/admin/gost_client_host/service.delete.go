package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/model"
	"server/repository"
	"server/service/common/cache"
	"server/service/gost_engine"
)

type DeleteReq struct {
	Code string `binding:"required" json:"code" label:"编号"`
}

func (service *service) Delete(req DeleteReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *gorm.DB) error {
		var host model.GostClientHost
		if tx.Preload("Node").Where("code = ?", req.Code).First(&host).RowsAffected == 0 {
			return nil
		}
		tx.Where("prefix = ? AND node_code = ?", host.DomainPrefix, host.NodeCode).Delete(&model.GostNodeDomain{})
		tx.Where("tunnel_code = ?", host.Code).Delete(&model.GostAuth{})
		if err := tx.Omit("Node").Delete(&host).Error; err != nil {
			log.Error("删除用户域名解析失败", zap.Error(err))
			return errors.New("操作失败")
		}
		gost_engine.ClientRemoveHostConfig(tx, host, host.Node)
		cache.DelTunnelInfo(req.Code)
		return nil
	})
}
