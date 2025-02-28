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
	"server/service/common/node_port"
	"server/service/common/node_rule"
	"server/service/gost_engine"
	"time"
)

type CreateReq struct {
	Name       string `binding:"required" json:"name" label:"名称"`
	TargetIp   string `binding:"required" json:"targetIp" label:"内网IP"`
	TargetPort string `binding:"required" json:"targetPort" label:"内网端口"`
	NoDelay    int    `json:"noDelay" label:"兼容模式"`
	ClientCode string `binding:"required" json:"clientCode" label:"客户端编号"`
	NodeCode   string `binding:"required" json:"nodeCode" label:"节点编号"`
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

	return db.Transaction(func(tx *gorm.DB) error {
		var user model.SystemUser
		if tx.Where("code = ?", claims.Code).First(&user).RowsAffected == 0 {
			return errors.New("用户错误")
		}

		var node model.GostNode
		if tx.Where("code = ?", req.NodeCode).First(&node).RowsAffected == 0 {
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

		var cfg model.GostNodeConfig
		if tx.Where("code = ? AND node_code = ?", req.ConfigCode, node.Code).First(&cfg).RowsAffected == 0 {
			return errors.New("套餐错误")
		}

		var expAt = time.Now().Unix()
		switch cfg.ChargingType {
		case model.GOST_CONFIG_CHARGING_CUCLE_DAY:
			expAt = time.Now().Add(time.Duration(cfg.Cycle) * 24 * time.Hour).Unix()
			if user.Amount.LessThan(cfg.Amount) {
				return errors.New("积分不足")
			}
			user.Amount = user.Amount.Sub(cfg.Amount)
			if err := tx.Save(&user).Error; err != nil {
				log.Error("扣减积分失败", zap.Error(err))
				return errors.New("操作失败")
			}
		case model.GOST_CONFIG_CHARGING_ONLY_ONCE:
			if user.Amount.LessThan(cfg.Amount) {
				return errors.New("积分不足")
			}
			user.Amount = user.Amount.Sub(cfg.Amount)
			if err := tx.Save(&user).Error; err != nil {
				log.Error("扣减积分失败", zap.Error(err))
				return errors.New("操作失败")
			}
		}

		var client model.GostClient
		if tx.Where("code = ? AND user_code = ?", req.ClientCode, claims.Code).First(&client).RowsAffected == 0 {
			return errors.New("客户端错误")
		}

		port, err := node_port.GetPort(node.Code)
		if err != nil {
			return err
		}
		if err = tx.Create(&model.GostNodePort{
			Port:     port,
			NodeCode: node.Code,
		}).Error; err != nil {
			node_port.ReleasePort(node.Code, port)
			log.Error("端口转发，端口冲突", zap.Error(err))
			return errors.New("操作失败")
		}

		var forward = model.GostClientForward{
			Name:       req.Name,
			TargetIp:   req.TargetIp,
			TargetPort: req.TargetPort,
			Port:       port,
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
		if err = tx.Create(&forward).Error; err != nil {
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
		if err = tx.Create(&auth).Error; err != nil {
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
