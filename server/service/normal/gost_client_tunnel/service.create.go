package service

import (
	"errors"
	"github.com/google/uuid"
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
	TargetIp   string `binding:"required" json:"targetIp" label:"内网IP"`
	TargetPort string `binding:"required" json:"targetPort" label:"内网端口"`
	NoDelay    int    `json:"noDelay" label:"兼容模式"`
	NodeCode   string `binding:"required" json:"nodeCode" label:"节点编号"`
	ClientCode string `binding:"required" json:"clientCode" label:"客户端编号"`
	ConfigCode string `binding:"required" json:"configCode" label:"套餐配置"`
}

func (service *service) Create(claims jwt.Claims, req CreateReq) error {
	db, _, log := repository.Get("")
	if !utils.ValidateLocalIP(req.TargetIp) {
		return errors.New("内网IP格式错误")
	}
	if !utils.ValidatePort(req.TargetPort) {
		return errors.New("内网端口格式错误")
	}

	var cfg model.SystemConfigGost
	cache.GetSystemConfigGost(&cfg)
	if cfg.FuncTunnel != "1" {
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
			user.Amount = user.Amount.Sub(cfg.Amount)
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

		var tunnel = model.GostClientTunnel{
			Name:       req.Name,
			TargetIp:   req.TargetIp,
			TargetPort: req.TargetPort,
			VKey:       uuid.NewString(),
			NoDelay:    req.NoDelay,
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
				OnlyChina:    cfg.OnlyChina,
				ExpAt:        expAt,
			},
		}
		if err := tx.GostClientTunnel.Create(&tunnel); err != nil {
			log.Error("新增用户私有隧道失败", zap.Error(err))
			return errors.New("操作失败")
		}
		var auth = model.GostAuth{
			TunnelType: model.GOST_TUNNEL_TYPE_TUNNEL,
			TunnelCode: tunnel.Code,
			User:       utils.RandStr(10, utils.AllDict),
			Password:   utils.RandStr(10, utils.AllDict),
		}
		if err := tx.GostAuth.Create(&auth); err != nil {
			log.Error("生成授权信息失败", zap.Error(err))
			return errors.New("操作失败")
		}
		cache.SetGostAuth(auth.User, auth.Password, tunnel.Code)
		gost_engine.ClientTunnelConfig(tx, tunnel.Code)
		cache.SetTunnelInfo(cache.TunnelInfo{
			Code:        tunnel.Code,
			Type:        model.GOST_TUNNEL_TYPE_TUNNEL,
			ClientCode:  tunnel.ClientCode,
			UserCode:    tunnel.UserCode,
			NodeCode:    tunnel.NodeCode,
			ChargingTye: tunnel.ChargingType,
			ExpAt:       tunnel.ExpAt,
			Limiter:     tunnel.Limiter,
		})
		return nil
	})
}
