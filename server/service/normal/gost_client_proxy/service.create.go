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
	"server/service/common/node_rule"
	"server/service/gost_engine"
	"time"
)

type CreateReq struct {
	Name       string `binding:"required" json:"name" label:"名称"`
	Port       string `binding:"required" json:"port" label:"本地端口"`
	Protocol   string `binding:"required" json:"protocol" label:"协议"`
	ClientCode string `binding:"required" json:"clientCode" label:"客户端编号"`
	NodeCode   string `binding:"required" json:"nodeCode" label:"节点编号"`
	ConfigCode string `binding:"required" json:"configCode" label:"套餐配置"`
}

func (service *service) Create(claims jwt.Claims, req CreateReq) error {
	db, _, log := repository.Get("")
	if !utils.ValidatePort(req.Port) {
		return errors.New("本地端口格式错误")
	}

	var cfg model.SystemConfigGost
	cache.GetSystemConfigGost(&cfg)
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
			if _, err := tx.SystemUser.Where(
				tx.SystemUser.Code.Eq(user.Code),
				tx.SystemUser.Version.Eq(user.Version),
			).UpdateSimple(
				tx.SystemUser.Amount.Value(user.Amount.Sub(user.Amount.Sub(cfg.Amount))),
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
				tx.SystemUser.Amount.Value(user.Amount.Sub(user.Amount.Sub(cfg.Amount))),
				tx.SystemUser.Version.Value(user.Version+1),
			); err != nil {
				log.Error("扣减积分失败", zap.Error(err))
				return errors.New("操作失败")
			}
		}

		if err := CheckPort(tx, *client, req.Port); err != nil {
			return err
		}

		var proxy = model.GostClientProxy{
			Name:       req.Name,
			Protocol:   req.Protocol,
			Port:       req.Port,
			NodeCode:   req.NodeCode,
			ClientCode: req.ClientCode,
			UserCode:   claims.Code,
			GostClientConfig: model.GostClientConfig{
				ChargingType: cfg.ChargingType,
				Cycle:        cfg.Cycle,
				Amount:       cfg.Amount,
				Limiter:      cfg.Limiter,
				RLimiter:     cfg.RLimiter,
				CLimiter:     cfg.CLimiter,
				ExpAt:        expAt,
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
		cache.SetGostAuth(auth.User, auth.Password, proxy.Code)
		gost_engine.ClientProxyConfig(tx, proxy.Code)
		cache.SetTunnelInfo(cache.TunnelInfo{
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

const (
	checkInterval = 200 * time.Millisecond
	maxRetries    = 25 // 5 seconds total (200ms * 25)
)

func CheckPort(tx *query.Query, client model.GostClient, port string) error {
	// 判断客户端状态，离线则不检测端口
	if !cache.GetClientOnline(client.Code) {
		return nil
	}
	// 发送端口检测消息
	async, allow := gost_engine.ClientPortCheck(tx, client.Code, port)
	if async {
		if allow {
			return nil
		} else {
			return errors.New("端口已被占用")
		}
	}
	for retry := 0; retry < maxRetries; retry++ {
		time.Sleep(checkInterval)
		if use, ok := cache.GetClientPortUse(client.Code, port); ok {
			if !use {
				return nil
			}
			return errors.New("端口已被占用")
		}
	}
	return errors.New("验证端口超时")
}
