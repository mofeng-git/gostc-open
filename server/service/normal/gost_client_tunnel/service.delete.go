package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/model"
	"server/pkg/jwt"
	"server/repository"
	"server/service/common/cache"
	"server/service/gost_engine"
)

type DeleteReq struct {
	Code string `binding:"required" json:"code" label:"编号"`
}

func (service *service) Delete(claims jwt.Claims, req DeleteReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *gorm.DB) error {
		var tunnel model.GostClientTunnel
		if tx.Where("code = ? AND user_code = ?", req.Code, claims.Code).First(&tunnel).RowsAffected == 0 {
			return nil
		}
		tx.Where("tunnel_code = ?", tunnel.Code).Delete(&model.GostAuth{})
		if err := tx.Delete(&tunnel).Error; err != nil {
			log.Error("删除用户私有隧道失败", zap.Error(err))
			return errors.New("操作失败")
		}
		gost_engine.ClientRemoveTunnelConfig(tx, tunnel, tunnel.Node)
		cache.DelTunnelInfo(req.Code)
		return nil
	})
}
