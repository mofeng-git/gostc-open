package node_rule

import "gorm.io/gorm"

func init() {
}

type DefaultRule struct{}

func (u DefaultRule) Name() string {
	return "无使用条件约束"
}

func (u DefaultRule) Code() string {
	return ""
}

func (u DefaultRule) Description() string {
	return "无使用条件约束"
}

func (u DefaultRule) Allow(tx *gorm.DB, userCode string) bool {
	return true
}
