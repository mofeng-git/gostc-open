package model

import "github.com/shopspring/decimal"

const (
	SYSTEM_IS_ADMIN = 1 // 是管理员
	SYSTEM_NO_ADMIN = 2 // 不是管理员
)

type SystemUser struct {
	Base
	Account  string          `gorm:"column:account;uniqueIndex;comment:账号"`
	Password string          `gorm:"column:password;comment:密码"`
	Salt     string          `gorm:"column:salt;size:8;comment:盐"`
	OtpKey   string          `gorm:"column:otp_key;index;default:'';not null"`
	Admin    int             `gorm:"column:admin;size:1;default:2;comment:是否为管理员"`
	Amount   decimal.Decimal `gorm:"column:amount;default:0;comment:余额"`
}
