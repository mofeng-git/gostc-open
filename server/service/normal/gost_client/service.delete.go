package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/model"
	"server/pkg/jwt"
	"server/repository"
	"server/service/gost_engine"
)

type DeleteReq struct {
	Code string `binding:"required" json:"code" label:"编号"`
}

func (service *service) Delete(claims jwt.Claims, req DeleteReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *gorm.DB) error {
		var client model.GostClient
		if tx.Where("code = ? AND user_code = ?", req.Code, claims.Code).First(&client).RowsAffected == 0 {
			return nil
		}
		var total int64
		tx.Model(&model.GostClientHost{}).Where("client_code = ?", client.Code).Count(&total)
		if total > 0 {
			return errors.New("请先删除该客户端的所有隧道")
		}
		tx.Model(&model.GostClientForward{}).Where("client_code = ?", client.Code).Count(&total)
		if total > 0 {
			return errors.New("请先删除该客户端的所有隧道")
		}
		tx.Model(&model.GostClientTunnel{}).Where("client_code = ?", client.Code).Count(&total)
		if total > 0 {
			return errors.New("请先删除该客户端的所有隧道")
		}
		if err := tx.Delete(&client).Error; err != nil {
			log.Error("删除客户端失败", zap.Error(err))
			return errors.New("操作失败")
		}
		tx.Where("client_code = ?", client.Code).Delete(&model.GostClientLogger{})
		gost_engine.ClientStop(client.Code, "客户端已被删除")
		return nil
	})
}
