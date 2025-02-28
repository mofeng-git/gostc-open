package node_rule

import (
	"gorm.io/gorm"
)

func init() {
}

type UserLevelRule struct{}

func (u UserLevelRule) Code() string {
	return "user_level_rule"
}

func (u UserLevelRule) Name() string {
	return "账号等级限制1级"
}

func (u UserLevelRule) Description() string {
	return "需要账号等级达到1级"
}

func (u UserLevelRule) Allow(tx *gorm.DB, userCode string) bool {
	// 允许使用
	return true
}
