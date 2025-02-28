package model

import "github.com/shopspring/decimal"

type SystemUserCheckin struct {
	Base
	UserCode  string          `gorm:"column:user_code;uniqueIndex:system_user_checkin_uidx"`
	Account   string          `gorm:"column:account;uniqueIndex:system_user_checkin_uidx"`
	EventDate string          `gorm:"column:event_date;uniqueIndex:system_user_checkin_uidx"`
	Amount    decimal.Decimal `gorm:"column:amount;index"`
}
