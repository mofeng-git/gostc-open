package service

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"server/model"
	"server/repository"
)

type CreateReq struct {
	Name         string `binding:"required" json:"name"`
	ChargingType int    `binding:"required" json:"chargingType" label:"计费类型"`
	Cycle        int    `json:"cycle" label:"计费周期(天)"`
	Amount       string `json:"amount" label:"计费金额"`
	Limiter      int    `json:"limiter" label:"速率"`
	RLimiter     int    `json:"rLimiter" label:"并发数"`
	CLimiter     int    `json:"cLimiter" label:"连接数"`
	OnlyChina    int    `json:"onlyChina" label:"仅中国大陆可用"`
	NodeCode     string `binding:"required" json:"nodeCode"`
	IndexValue   int    `json:"indexValue"`
}

func (service *service) Create(req CreateReq) (err error) {
	db, _, log := repository.Get("")
	var amount decimal.Decimal
	switch req.ChargingType {
	case model.GOST_CONFIG_CHARGING_ONLY_ONCE, model.GOST_CONFIG_CHARGING_CUCLE_DAY:
		amount, err = decimal.NewFromString(req.Amount)
		if err != nil {
			return fmt.Errorf("金额错误，%v", err)
		}
		if req.ChargingType == model.GOST_CONFIG_CHARGING_CUCLE_DAY && req.Cycle <= 0 {
			return errors.New("计费循环周期错误")
		}
	}

	if err = db.GostNodeConfig.Create(&model.GostNodeConfig{
		Name:         req.Name,
		ChargingType: req.ChargingType,
		Cycle:        req.Cycle,
		Amount:       amount,
		Limiter:      req.Limiter,
		RLimiter:     req.RLimiter,
		CLimiter:     req.CLimiter,
		NodeCode:     req.NodeCode,
		OnlyChina:    req.OnlyChina,
		IndexValue:   req.IndexValue,
	}); err != nil {
		log.Error("新增套餐配置失败", zap.Error(err))
		return errors.New("操作失败")
	}
	return nil
}
