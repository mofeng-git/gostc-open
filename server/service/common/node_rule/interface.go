package node_rule

import (
	"gorm.io/gorm"
)

type RuleInterface interface {
	Code() string
	Allow(db *gorm.DB, userCode string) bool
	Name() string
	Description() string
}

var RuleList []RuleInterface
var RuleMap = make(map[string]RuleInterface)

func init() {
	RuleList = append(RuleList,
		DefaultRule{},
		UserLevelRule{},
		UserQQGroupRule{})
	for _, rule := range RuleList {
		RuleMap[rule.Code()] = rule
	}
}
