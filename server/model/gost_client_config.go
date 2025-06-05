package model

import (
	"github.com/shopspring/decimal"
)

type GostClientConfig struct {
	ChargingType int             `gorm:"column:charging_type;default:1;index;comment:计费方式"`
	Cycle        int             `gorm:"column:cycle;default:0;comment:计费周期(天)"`
	Amount       decimal.Decimal `gorm:"column:amount;default:0;comment:费用"`
	Limiter      int             `gorm:"column:limiter;comment:速度速率限制"`
	RLimiter     int             `gorm:"column:r_limiter;comment:并发数量限制"`
	CLimiter     int             `gorm:"column:c_limiter;comment:连接数量限制"`
	ExpAt        int64           `gorm:"column:exp_at;index;comment:套餐到期时间"`
}
