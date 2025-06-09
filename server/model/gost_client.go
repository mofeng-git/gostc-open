package model

type GostClient struct {
	Base
	Name     string     `gorm:"column:name;index;comment:名称"`
	UserCode string     `gorm:"column:user_code;index;comment:用户编号"`
	User     SystemUser `gorm:"foreignKey:UserCode;references:Code"`
	Key      string     `gorm:"column:key;size:100;uniqueIndex;comment:连接密钥"`
}
