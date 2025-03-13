package service

import (
	"errors"
	"go.uber.org/zap"
	"server/model"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/repository/query"
	"server/service/common/cache"
	"server/service/common/node_port"
	"server/service/common/node_rule"
	"server/service/gost_engine"
	"time"
)

type CreateReq struct {
	Name          string `binding:"required" json:"name" label:"名称"`
	TargetIp      string `binding:"required" json:"targetIp" label:"内网IP"`
	TargetPort    string `binding:"required" json:"targetPort" label:"内网端口"`
	ProxyProtocol int    `json:"proxyProtocol"`
	ClientCode    string `binding:"required" json:"clientCode" label:"客户端编号"`
	NodeCode      string `binding:"required" json:"nodeCode" label:"节点编号"`
	ConfigCode    string `binding:"required" json:"configCode" label:"套餐配置"`
}

func GetPort(tx *query.Query, node model.GostNode) (port string, err error) {
	for {
		port, err = node_port.GetPort(node.Code)
		if err != nil {
			return "", err
		}
		version := cache.GetNodeVersion(node.Code)
		if version >= "v1.1.7" {
			gost_engine.NodePortCheck(tx, node.Code, port)
			var available bool // 是否可用
			var retry int
			for {
				time.Sleep(time.Millisecond * 200)
				retry++
				use, ok := cache.GetNodePortUse(node.Code, port)
				if ok {
					if use {
						available = false
					} else {
						available = true
					}
					break
				}
				if retry > 5*5 {
					break
				}
			}
			// 可用，则结束获取端口
			if available {
				break
			}
		} else {
			return port, err
		}
	}
	return port, nil
}

func (service *service) Create(claims jwt.Claims, req CreateReq) error {
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

		node, _ := tx.GostNode.Where(tx.GostNode.Code.Eq(req.NodeCode)).First()
		if node == nil {
			return errors.New("节点错误")
		}
		if node.Forward != 1 {
			return errors.New("该节点未启用端口转发功能")
		}

		for _, ruleCode := range node.GetRules() {
			rule := node_rule.RuleMap[ruleCode]
			if rule.Code() == "" {
				continue
			}
			if !rule.Allow(tx, user.Code) {
				return errors.New("规则不符合，" + rule.Description())
			}
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
			user.Amount = user.Amount.Sub(cfg.Amount)
			if err := tx.SystemUser.Save(user); err != nil {
				log.Error("扣减积分失败", zap.Error(err))
				return errors.New("操作失败")
			}
		case model.GOST_CONFIG_CHARGING_ONLY_ONCE:
			if user.Amount.LessThan(cfg.Amount) {
				return errors.New("积分不足")
			}
			user.Amount = user.Amount.Sub(cfg.Amount)
			if err := tx.SystemUser.Save(user); err != nil {
				log.Error("扣减积分失败", zap.Error(err))
				return errors.New("操作失败")
			}
		}

		port, err := GetPort(tx, *node)
		if err != nil {
			return err
		}

		if err = tx.GostNodePort.Create(&model.GostNodePort{
			Port:     port,
			NodeCode: node.Code,
		}); err != nil {
			node_port.ReleasePort(node.Code, port)
			log.Error("端口转发，端口冲突", zap.Error(err))
			return errors.New("操作失败")
		}

		var forward = model.GostClientForward{
			Name:          req.Name,
			TargetIp:      req.TargetIp,
			TargetPort:    req.TargetPort,
			Port:          port,
			ProxyProtocol: req.ProxyProtocol,
			NodeCode:      req.NodeCode,
			ClientCode:    req.ClientCode,
			UserCode:      claims.Code,
			GostClientConfig: model.GostClientConfig{
				ChargingType: cfg.ChargingType,
				Cycle:        cfg.Cycle,
				Amount:       cfg.Amount,
				Limiter:      cfg.Limiter,
				RLimiter:     cfg.RLimiter,
				CLimiter:     cfg.CLimiter,
				OnlyChina:    cfg.OnlyChina,
				ExpAt:        expAt,
			},
		}
		if err = tx.GostClientForward.Create(&forward); err != nil {
			node_port.ReleasePort(node.Code, port)
			log.Error("新增用户端口转发失败", zap.Error(err))
			return errors.New("操作失败")
		}
		var auth = model.GostAuth{
			TunnelType: model.GOST_TUNNEL_TYPE_FORWARD,
			TunnelCode: forward.Code,
			User:       utils.RandStr(10, utils.AllDict),
			Password:   utils.RandStr(10, utils.AllDict),
		}
		if err = tx.GostAuth.Create(&auth); err != nil {
			node_port.ReleasePort(node.Code, port)
			log.Error("生成授权信息失败", zap.Error(err))
			return errors.New("操作失败")
		}
		cache.SetGostAuth(auth.User, auth.Password, forward.Code)
		gost_engine.ClientForwardConfig(tx, forward.Code)
		cache.SetTunnelInfo(cache.TunnelInfo{
			Code:        forward.Code,
			Type:        model.GOST_TUNNEL_TYPE_FORWARD,
			ClientCode:  forward.ClientCode,
			UserCode:    forward.UserCode,
			NodeCode:    forward.NodeCode,
			ChargingTye: forward.ChargingType,
			ExpAt:       forward.ExpAt,
			Limiter:     forward.Limiter,
		})
		return nil
	})
}
