package model

type SystemUserEmail struct {
	Base
	Email    string `gorm:"column:email;size:200;uniqueIndex;comment:邮箱"`
	UserCode string `gorm:"column:user_code;size:100;uniqueIndex;comment:用户编号"`
}
