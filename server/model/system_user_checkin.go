package model

import "github.com/shopspring/decimal"

type SystemUserCheckin struct {
	Base
	UserCode  string          `gorm:"column:user_code;size:100;uniqueIndex:system_user_checkin_uidx"`
	Account   string          `gorm:"column:account;size:100;uniqueIndex:system_user_checkin_uidx"`
	EventDate string          `gorm:"column:event_date;size:20;uniqueIndex:system_user_checkin_uidx"`
	Amount    decimal.Decimal `gorm:"column:amount;index"`
}
