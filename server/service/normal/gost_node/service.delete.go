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
		var nodeBind model.GostNodeBind
		if tx.Where("node_code = ? AND user_code = ?", req.Code, claims.Code).First(&nodeBind).RowsAffected == 0 {
			return nil
		}

		var node model.GostNode
		if tx.Where("code = ?", req.Code).First(&node).RowsAffected == 0 {
			return nil
		}

		var hostTotal int64
		tx.Model(&model.GostClientHost{}).Where("node_code = ?", node.Code).Count(&hostTotal)
		if hostTotal > 0 {
			return errors.New("请先删除该节点的所有隧道")
		}
		var forwardTotal int64
		tx.Model(&model.GostClientForward{}).Where("node_code = ?", node.Code).Count(&forwardTotal)
		if forwardTotal > 0 {
			return errors.New("请先删除该节点的所有隧道")
		}
		var tunnelTotal int64
		tx.Model(&model.GostClientTunnel{}).Where("node_code = ?", node.Code).Count(&tunnelTotal)
		if tunnelTotal > 0 {
			return errors.New("请先删除该节点的所有隧道")
		}
		if err := tx.Delete(&node).Error; err != nil {
			log.Error("删除节点失败", zap.Error(err))
			return errors.New("操作失败")
		}
		tx.Where("node_code = ?", node.Code).Delete(&model.GostNodeConfig{})
		tx.Where("node_code = ?", node.Code).Delete(&model.GostNodeLogger{})
		tx.Where("node_code = ?", node.Code).Delete(&model.GostNodeBind{})
		gost_engine.NodeStop(node.Code, "节点已被删除")
		return nil
	})
}
