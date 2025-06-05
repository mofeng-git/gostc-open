package model

import (
	"github.com/shopspring/decimal"
)

const (
	GOST_CONFIG_CHARGING_ONLY_ONCE = 1 // 一次性计费
	GOST_CONFIG_CHARGING_CUCLE_DAY = 2 // 周期计费(天)
	GOST_CONFIG_CHARGING_FREE      = 3 // 免费

)

type GostNodeConfig struct {
	Base
	IndexValue   int             `gorm:"column:index_value;index;default:1000;comment:排序，升序"`
	Name         string          `gorm:"column:name;index;comment:名称"`
	ChargingType int             `gorm:"column:charging_type;default:1;index;comment:计费方式"`
	Cycle        int             `gorm:"column:cycle;default:0;comment:计费周期(天)"`
	Amount       decimal.Decimal `gorm:"column:amount;default:0;comment:费用"`
	Limiter      int             `gorm:"column:limiter;comment:速率(mbps)"`
	RLimiter     int             `gorm:"column:r_limiter;comment:并发数"`
	CLimiter     int             `gorm:"column:c_limiter;comment:连接数"`
	NodeCode     string          `gorm:"column:node_code;index;comment:节点编号"`
	Node         GostNode        `gorm:"foreignKey:NodeCode;references:Code"`
}
