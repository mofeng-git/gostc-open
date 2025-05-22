package service

import (
	"errors"
	"go.uber.org/zap"
	"server/global"
	"server/model"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/repository/query"
	"server/service/common/cache"
	"server/service/common/node_port"
	"server/service/gost_engine"
)

type UpdateReq struct {
	Code          string `binding:"required" json:"code"`
	Name          string `binding:"required" json:"name"`
	TargetIp      string `binding:"required" json:"targetIp"`
	TargetPort    string `binding:"required" json:"targetPort"`
	Port          string `json:"port"`
	ProxyProtocol int    `json:"proxyProtocol"`
}

func (service *service) Update(claims jwt.Claims, req UpdateReq) error {
	db, _, log := repository.Get("")
	if !utils.ValidateLocalIP(req.TargetIp) {
		return errors.New("内网IP格式错误")
	}
	if !utils.ValidatePort(req.TargetPort) {
		return errors.New("内网端口格式错误")
	}

	if req.Port != "" && !utils.ValidatePort(req.TargetPort) {
		return errors.New("远程端口格式错误")
	}
	return db.Transaction(func(tx *query.Query) error {
		user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(claims.Code)).First()
		if user == nil {
			return errors.New("用户错误")
		}

		forward, _ := tx.GostClientForward.Where(
			tx.GostClientForward.UserCode.Eq(user.Code),
			tx.GostClientForward.Code.Eq(req.Code),
		).First()
		if forward == nil {
			return errors.New("操作失败")
		}

		forward.Name = req.Name
		forward.TargetIp = req.TargetIp
		forward.TargetPort = req.TargetPort
		forward.ProxyProtocol = req.ProxyProtocol

		var oldPort = forward.Port
		if req.Port != "" && req.Port != oldPort {
			forward.Port = req.Port
			if !node_port.ValidPort(forward.NodeCode, req.Port, true) {
				return errors.New("新的端口未开放或已被占用")
			}

			if cache.GetNodeOnline(forward.NodeCode) {
				if !validPortAvailable(tx, forward.NodeCode, req.Port) {
					return errors.New("新的端口已被占用")
				}
			}

			// 删除旧端口资源
			if _, err := tx.GostNodePort.Where(
				tx.GostNodePort.Port.Eq(oldPort),
				tx.GostNodePort.NodeCode.Eq(forward.NodeCode),
			).Delete(); err != nil {
				global.Logger.Error("端口转发，修改远程端口，释放旧端口资源失败", zap.Error(err))
				return errors.New("释放旧端口资源失败")
			}

			// 分配新端口资源
			if err := tx.GostNodePort.Create(&model.GostNodePort{
				Port:     req.Port,
				NodeCode: forward.NodeCode,
			}); err != nil {
				return errors.New("新的端口已被占用")
			}
		}

		if err := tx.GostClientForward.Save(forward); err != nil {
			log.Error("修改端口转发失败", zap.Error(err))
			return errors.New("操作失败")
		}
		if req.Port != "" && oldPort != req.Port {
			// 归还旧端口资源
			node_port.ReleasePort(forward.NodeCode, oldPort)
		}
		gost_engine.ClientForwardConfig(tx, forward.Code)
		return nil
	})
}
