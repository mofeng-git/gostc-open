package node_rule

import (
	"server/repository/query"
	"sync"
)

type RuleInterface interface {
	Code() string
	Allow(db *query.Query, userCode string) bool
	Name() string
	Description() string
}

var Registry = registry{
	sort: []string{defaultRule.Code()},
	data: map[string]RuleInterface{
		defaultRule.Code(): defaultRule,
	},
	mu: &sync.RWMutex{},
}

type registry struct {
	sort []string
	data map[string]RuleInterface
	mu   *sync.RWMutex
}

func (r *registry) GetRules() (result []RuleInterface) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, code := range r.sort {
		result = append(result, r.data[code])
	}
	return result
}
func (r *registry) GetRule(code string) RuleInterface {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rule := r.data[code]
	if rule == nil {
		return defaultRule
	}
	return rule
}

func (r *registry) DelRule(code string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var newSort []string
	for _, item := range r.sort {
		if item != code {
			newSort = append(newSort, item)
		}
	}
	r.sort = newSort
	delete(r.data, code)
}

func (r *registry) SetRule(rule RuleInterface) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[rule.Code()]; ok {
		r.data[rule.Code()] = rule
	} else {
		r.sort = append(r.sort, rule.Code())
		r.data[rule.Code()] = rule
	}
}
