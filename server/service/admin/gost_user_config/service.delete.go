package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/model"
	"server/repository"
)

type DeleteReq struct {
	Code string `binding:"required" json:"code" label:"编号"`
}

func (service *service) Delete(req DeleteReq) error {
	db, _, log := repository.Get("")
	var cfg model.GostClientConfig
	if db.Where("code = ?", req.Code).First(&cfg).RowsAffected == 0 {
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&cfg).Error; err != nil {
			log.Error("删除用户套餐配置失败", zap.Error(err))
			return errors.New("删除失败")
		}
		//if cfg.TunnelCode != "" {
		//	switch cfg.TunnelType {
		//	case model.GOST_TUNNEL_TYPE_HOST:
		//		var host model.GostClientHost
		//		if tx.Where("code = ?", cfg.TunnelCode).First(&host).RowsAffected != 0 {
		//			gost_engine.ClientRemoveConfig(host.ClientCode, []string{
		//				host.Code,
		//			})
		//		}
		//		tx.Where("domain_prefix = ? AND node_code = ?", host.DomainPrefix, host.NodeCode).Delete(&model.GostNodeDomain{})
		//		tx.Delete(&host)
		//	case model.GOST_TUNNEL_TYPE_FORWARD:
		//		var forward model.GostClientForward
		//		if tx.Where("code = ?", cfg.TunnelCode).First(&forward).RowsAffected != 0 {
		//			gost_engine.ClientRemoveConfig(forward.ClientCode, []string{
		//				"tcp_" + forward.Code,
		//				"udp_" + forward.Code,
		//			})
		//		}
		//		node_port.ReleasePort(forward.NodeCode, forward.Port)
		//		tx.Delete(&forward)
		//	case model.GOST_TUNNEL_TYPE_TUNNEL:
		//		var tunnel model.GostClientTunnel
		//		if tx.Where("code = ?", cfg.TunnelCode).First(&tunnel).RowsAffected != 0 {
		//			gost_engine.ClientRemoveConfig(tunnel.ClientCode, []string{
		//				"tcp_" + tunnel.Code,
		//				"udp_" + tunnel.Code,
		//			})
		//		}
		//		tx.Delete(&tunnel)
		//	}
		//	tx.Where("tunnel_code = ?", cfg.TunnelCode).Delete(&model.GostAuth{})
		//}

		return nil
	})
}
