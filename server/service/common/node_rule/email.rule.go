package node_rule

import (
	"server/repository/query"
)

func init() {
	Registry.SetRule(emailRule)
}

var emailRule = &EmailRule{}

type EmailRule struct{}

func (u EmailRule) Name() string {
	return "强制绑定邮箱"
}

func (u EmailRule) Code() string {
	return "bind_email"
}
func (u EmailRule) Group() string {
	return "bind_email"
}

func (u EmailRule) Description() string {
	return "必须绑定邮箱才能使用"
}

func (u EmailRule) Allow(tx *query.Query, userCode string) bool {
	userEmail, err := tx.SystemUserEmail.Where(tx.SystemUserEmail.UserCode.Eq(userCode)).First()
	if err != nil {
		return false
	}
	if userEmail.Email == "" {
		return false
	}
	return true
}
