package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/model"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/service/common/cache"
	"server/service/gost_engine"
)

type UpdateReq struct {
	Code         string `binding:"required" json:"code"`
	Name         string `binding:"required" json:"name"`
	TargetIp     string `binding:"required" json:"targetIp"`
	TargetPort   string `binding:"required" json:"targetPort"`
	DomainPrefix string `binding:"required" json:"domainPrefix"`
}

func (service *service) Update(claims jwt.Claims, req UpdateReq) error {
	db, _, log := repository.Get("")
	if !utils.ValidateLocalIP(req.TargetIp) {
		return errors.New("内网IP格式错误")
	}
	if !utils.ValidatePort(req.TargetPort) {
		return errors.New("内网端口格式错误")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		var user model.SystemUser
		if tx.Where("code = ?", claims.Code).First(&user).RowsAffected == 0 {
			return errors.New("用户错误")
		}

		var host model.GostClientHost
		if tx.Preload("Node").Where("code = ? AND user_code = ?", req.Code, user.Code).First(&host).RowsAffected == 0 {
			return errors.New("操作失败")
		}

		if req.DomainPrefix != host.DomainPrefix {
			if err := tx.Create(&model.GostNodeDomain{
				Prefix:   req.DomainPrefix,
				NodeCode: host.NodeCode,
			}).Error; err != nil {
				return errors.New("该域名前缀已被使用")
			}
			tx.Where("prefix = ? AND node_code = ?", host.DomainPrefix, host.NodeCode).Delete(&model.GostNodeDomain{})
			// 清除旧缓存
			cache.DelIngress(host.DomainPrefix + "." + host.Node.Domain)
			cache.DelIngress(host.DomainPrefix + "." + host.Node.Domain + ":" + host.Node.TunnelInPort)
		}

		host.Name = req.Name
		host.TargetIp = req.TargetIp
		host.TargetPort = req.TargetPort
		host.DomainPrefix = req.DomainPrefix

		if err := tx.Save(&host).Error; err != nil {
			log.Error("修改域名解析失败", zap.Error(err))
			return errors.New("操作失败")
		}

		gost_engine.ClientHostConfig(tx, host.Code)
		return nil
	})
}
