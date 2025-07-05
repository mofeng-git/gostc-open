package service

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"server/model"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	cache2 "server/repository/cache"
	"server/repository/query"
	"server/service/common/node_port"
	"server/service/common/node_rule"
	"server/service/engine"
	"time"
)

type CreateReq struct {
	Name       string `binding:"required" json:"name" label:"名称"`
	Port       string `json:"port" label:"本地端口"`
	AuthUser   string `json:"authUser"`
	AuthPwd    string `json:"authPwd"`
	Protocol   string `binding:"required" json:"protocol" label:"协议"`
	ClientCode string `binding:"required" json:"clientCode" label:"客户端编号"`
	NodeCode   string `binding:"required" json:"nodeCode" label:"节点编号"`
	ConfigCode string `binding:"required" json:"configCode" label:"套餐配置"`
}

func (service *service) Create(claims jwt.Claims, req CreateReq) error {
	db, _, log := repository.Get("")
	if req.Port != "" && !utils.ValidatePort(req.Port) {
		return errors.New("本地端口格式错误")
	}

	var cfg model.SystemConfigGost
	cache2.GetSystemConfigGost(&cfg)
	if cfg.FuncProxy != "1" {
		return errors.New("管理员未启用该功能")
	}

	return db.Transaction(func(tx *query.Query) error {
		user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(claims.Code)).First()
		if user == nil {
			return errors.New("用户错误")
		}

		node, _ := tx.GostNode.Where(tx.GostNode.Code.Eq(req.NodeCode)).First()
		if node == nil {
			return errors.New("节点错误")
		}
		if node.Proxy != 1 {
			return errors.New("该节点未启代理隧道功能")
		}

		if err := node_rule.VerifyAll(tx, user.Code, node.GetRules()); err != nil {
			return err
		}

		cfg, _ := tx.GostNodeConfig.Where(
			tx.GostNodeConfig.Code.Eq(req.ConfigCode),
			tx.GostNodeConfig.NodeCode.Eq(node.Code),
		).First()
		if cfg == nil {
			return errors.New("套餐错误")
		}
		client, _ := tx.GostClient.Where(
			tx.GostClient.UserCode.Eq(claims.Code),
			tx.GostClient.Code.Eq(req.ClientCode),
		).First()
		if client == nil {
			return errors.New("客户端错误")
		}

		var expAt = time.Now().Unix()
		switch cfg.ChargingType {
		case model.GOST_CONFIG_CHARGING_CUCLE_DAY:
			expAt = time.Now().Add(time.Duration(cfg.Cycle) * 24 * time.Hour).Unix()
			if user.Amount.LessThan(cfg.Amount) {
				return errors.New("积分不足")
			}
			if _, err := tx.SystemUser.Where(
				tx.SystemUser.Code.Eq(user.Code),
				tx.SystemUser.Version.Eq(user.Version),
			).UpdateSimple(
				tx.SystemUser.Amount.Value(user.Amount.Sub(cfg.Amount)),
				tx.SystemUser.Version.Value(user.Version+1),
			); err != nil {
				log.Error("扣减积分失败", zap.Error(err))
				return errors.New("操作失败")
			}
		case model.GOST_CONFIG_CHARGING_ONLY_ONCE:
			if user.Amount.LessThan(cfg.Amount) {
				return errors.New("积分不足")
			}
			if _, err := tx.SystemUser.Where(
				tx.SystemUser.Code.Eq(user.Code),
				tx.SystemUser.Version.Eq(user.Version),
			).UpdateSimple(
				tx.SystemUser.Amount.Value(user.Amount.Sub(cfg.Amount)),
				tx.SystemUser.Version.Value(user.Version+1),
			); err != nil {
				log.Error("扣减积分失败", zap.Error(err))
				return errors.New("操作失败")
			}
		}

		var err error
		var port = req.Port
		if port == "" {
			port, err = GetPort(tx, *node)
			if err != nil {
				return err
			}
		} else {
			if !node_port.ValidPort(node.Code, req.Port, true) {
				return errors.New("端口未开放或已被占用")
			}
			if cache2.GetNodeOnline(req.NodeCode) {
				if !validPortAvailable(tx, node.Code, req.Port) {
					return errors.New("端口已被占用")
				}
			}
		}

		if err = tx.GostNodePort.Create(&model.GostNodePort{
			Port:     port,
			NodeCode: node.Code,
		}); err != nil {
			node_port.ReleasePort(node.Code, port)
			log.Error("端口转发，端口冲突", zap.Error(err))
			return errors.New("操作失败")
		}

		var proxy = model.GostClientProxy{
			Name:       req.Name,
			Protocol:   req.Protocol,
			Port:       port,
			AuthUser:   req.AuthUser,
			AuthPwd:    req.AuthPwd,
			NodeCode:   req.NodeCode,
			Node:       model.GostNode{},
			ClientCode: req.ClientCode,
			Client:     model.GostClient{},
			UserCode:   claims.Code,
			User:       model.SystemUser{},
			GostClientConfig: model.GostClientConfig{
				ChargingType: cfg.ChargingType,
				Cycle:        cfg.Cycle,
				Amount:       cfg.Amount,
				Limiter:      cfg.Limiter,
				//RLimiter:     cfg.RLimiter,
				//CLimiter:     cfg.CLimiter,
				ExpAt: expAt,
			},
		}
		if err := tx.GostClientProxy.Create(&proxy); err != nil {
			log.Error("新增用户代理隧道失败", zap.Error(err))
			return errors.New("操作失败")
		}
		var auth = model.GostAuth{
			TunnelType: model.GOST_TUNNEL_TYPE_PROXY,
			TunnelCode: proxy.Code,
			User:       utils.RandStr(10, utils.AllDict),
			Password:   utils.RandStr(10, utils.AllDict),
		}
		if err := tx.GostAuth.Create(&auth); err != nil {
			log.Error("生成授权信息失败", zap.Error(err))
			return errors.New("操作失败")
		}
		cache2.SetGostAuth(auth.User, auth.Password, proxy.Code)
		engine.ClientProxyConfig(tx, proxy.Code)
		cache2.SetTunnelInfo(cache2.TunnelInfo{
			Code:        proxy.Code,
			Type:        model.GOST_TUNNEL_TYPE_PROXY,
			ClientCode:  proxy.ClientCode,
			UserCode:    proxy.UserCode,
			NodeCode:    proxy.NodeCode,
			ChargingTye: proxy.ChargingType,
			ExpAt:       proxy.ExpAt,
			Limiter:     proxy.Limiter,
		})
		return nil
	})
}

func validPortAvailable(tx *query.Query, nodeCode string, port string) bool {
	return engine.NodePortCheck(tx, nodeCode, port) == nil
}
func GetPort(tx *query.Query, node model.GostNode) (string, error) {
	for {
		port, err := node_port.GetPort(node.Code)
		if err != nil {
			return "", fmt.Errorf("failed to get port: %w", err)
		}
		// 判断版本，对离线节点，不检查端口
		if !cache2.GetNodeOnline(node.Code) || cache2.GetNodeVersion(node.Code) < "v1.1.7" {
			return port, nil
		}
		// 检测端口
		if validPortAvailable(tx, node.Code, port) {
			return port, nil
		}
		// 如果端口不可用，则继续循环
	}
}
