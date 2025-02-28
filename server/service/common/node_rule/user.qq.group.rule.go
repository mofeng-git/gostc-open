package node_rule

import (
	"gorm.io/gorm"
)

type UserQQGroupRule struct{}

func (u UserQQGroupRule) Code() string {
	return "user_qq_group_rule"
}

func (u UserQQGroupRule) Name() string {
	return "绑定QQ号"
}

func (u UserQQGroupRule) Description() string {
	return "需要绑定QQ号"
}

func (u UserQQGroupRule) Allow(tx *gorm.DB, userCode string) bool {
	//var bindQQ model.SystemUserQQ
	//if err := tx.Where("user_code = ?", userCode).First(&bindQQ).Error; err != nil {
	//	return false
	//}
	//return bindQQ.QQ != ""
	return true
}
