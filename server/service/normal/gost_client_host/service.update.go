package service

import (
	"errors"
	"go.uber.org/zap"
	"server/model"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/repository/query"
	"server/service/engine"
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

	return db.Transaction(func(tx *query.Query) error {
		user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(claims.Code)).First()
		if user == nil {
			return errors.New("用户错误")
		}

		host, _ := tx.GostClientHost.Where(
			tx.GostClientHost.UserCode.Eq(user.Code),
			tx.GostClientHost.Code.Eq(req.Code),
		).First()
		if host == nil {
			return errors.New("操作失败")
		}
		if req.DomainPrefix != host.DomainPrefix {
			if err := tx.GostNodeDomain.Create(&model.GostNodeDomain{
				Prefix:   req.DomainPrefix,
				NodeCode: host.NodeCode,
			}); err != nil {
				return errors.New("该域名前缀已被使用")
			}
			_, _ = tx.GostNodeDomain.Where(tx.GostNodeDomain.Prefix.Eq(host.DomainPrefix), tx.GostNodeDomain.NodeCode.Eq(host.NodeCode)).Delete()
			//清除旧缓存
			//cache.DelIngress(host.DomainPrefix + "." + host.Node.Domain)
			//cache.DelIngress(host.DomainPrefix + "." + host.Node.Domain + ":" + host.Node.TunnelInPort)
		}

		host.Name = req.Name
		host.TargetIp = req.TargetIp
		host.TargetPort = req.TargetPort
		host.DomainPrefix = req.DomainPrefix

		if err := tx.GostClientHost.Save(host); err != nil {
			log.Error("修改域名解析失败", zap.Error(err))
			return errors.New("操作失败")
		}

		engine.ClientHostConfig(tx, host.Code)
		return nil
	})
}
