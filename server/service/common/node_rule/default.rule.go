package node_rule

import (
	"server/repository/query"
)

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

func (u DefaultRule) Allow(tx *query.Query, userCode string) bool {
	return true
}
