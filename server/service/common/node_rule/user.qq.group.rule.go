package node_rule

import (
	"server/repository/query"
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

func (u UserQQGroupRule) Allow(tx *query.Query, userCode string) bool {
	return true
}
