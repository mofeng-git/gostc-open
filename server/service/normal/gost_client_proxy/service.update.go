package service

import (
	"errors"
	"go.uber.org/zap"
	"server/global"
	"server/model"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/repository/cache"
	"server/repository/query"
	"server/service/common/node_port"
	"server/service/engine"
)

type UpdateReq struct {
	Code           string `binding:"required" json:"code"`
	Name           string `binding:"required" json:"name"`
	Port           string `json:"port" label:"本地端口"`
	Protocol       string `binding:"required" json:"protocol" label:"协议"`
	AuthUser       string `json:"authUser"`
	AuthPwd        string `json:"authPwd"`
	UseEncryption  int    `json:"useEncryption"`
	UseCompression int    `json:"useCompression"`
	PoolCount      int    `json:"poolCount"`
}

func (service *service) Update(claims jwt.Claims, req UpdateReq) error {
	db, _, log := repository.Get("")
	if req.Port != "" && !utils.ValidatePort(req.Port) {
		return errors.New("内网端口格式错误")
	}

	return db.Transaction(func(tx *query.Query) error {
		user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(claims.Code)).First()
		if user == nil {
			return errors.New("用户错误")
		}

		proxy, _ := tx.GostClientProxy.Where(
			tx.GostClientProxy.UserCode.Eq(user.Code),
			tx.GostClientProxy.Code.Eq(req.Code),
		).First()
		if proxy == nil {
			return errors.New("操作失败")
		}

		client, _ := tx.GostClient.Where(
			tx.GostClient.UserCode.Eq(claims.Code),
			tx.GostClient.Code.Eq(proxy.ClientCode),
		).First()
		if client == nil {
			return errors.New("客户端错误")
		}

		proxy.Name = req.Name
		proxy.Protocol = req.Protocol
		proxy.AuthUser = req.AuthUser
		proxy.AuthPwd = req.AuthPwd
		proxy.UseEncryption = req.UseEncryption
		proxy.UseCompression = req.UseCompression
		proxy.PoolCount = req.PoolCount

		var oldPort = proxy.Port
		if req.Port != "" && req.Port != oldPort {
			proxy.Port = req.Port
			if !node_port.ValidPort(proxy.NodeCode, req.Port, true) {
				return errors.New("新的端口未开放或已被占用")
			}

			if cache.GetNodeOnline(proxy.NodeCode) {
				if !validPortAvailable(tx, proxy.NodeCode, req.Port) {
					return errors.New("新的端口已被占用")
				}
			}

			// 删除旧端口资源
			if _, err := tx.GostNodePort.Where(
				tx.GostNodePort.Port.Eq(oldPort),
				tx.GostNodePort.NodeCode.Eq(proxy.NodeCode),
			).Delete(); err != nil {
				global.Logger.Error("端口转发，修改远程端口，释放旧端口资源失败", zap.Error(err))
				return errors.New("释放旧端口资源失败")
			}

			// 分配新端口资源
			if err := tx.GostNodePort.Create(&model.GostNodePort{
				Port:     req.Port,
				NodeCode: proxy.NodeCode,
			}); err != nil {
				return errors.New("新的端口已被占用")
			}
		}

		if err := tx.GostClientProxy.Save(proxy); err != nil {
			log.Error("修改代理隧道失败", zap.Error(err))
			return errors.New("操作失败")
		}
		engine.ClientProxyConfig(tx, proxy.Code)
		return nil
	})
}
