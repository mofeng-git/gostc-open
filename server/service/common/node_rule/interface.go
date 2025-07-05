package node_rule

import (
	"errors"
	"server/repository/query"
	"strings"
	"sync"
)

type RuleInterface interface {
	Code() string
	Allow(db *query.Query, userCode string) bool
	Group() string
	Name() string
	Description() string
}

var Registry = registry{
	sort: []string{},
	data: map[string]RuleInterface{},
	mu:   &sync.RWMutex{},
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

func VerifyAll(tx *query.Query, userCode string, codes []string) error {
	var groupMap = make(map[string][]RuleInterface)
	for _, code := range codes {
		if code == "" {
			continue
		}
		rule := Registry.GetRule(code)
		if rule.Code() == "" {
			continue
		}
		group := rule.Group()
		groupMap[group] = append(groupMap[group], rule)
	}
	for _, rules := range groupMap {
		if err := verifyGroup(tx, userCode, rules); err != nil {
			return errors.New("规则不符合，" + err.Error())
		}
	}
	return nil
}

func verifyGroup(tx *query.Query, userCode string, rules []RuleInterface) error {
	var errMsg []string
	for _, rule := range rules {
		if rule.Allow(tx, userCode) {
			errMsg = []string{}
			break
		} else {
			errMsg = append(errMsg, rule.Description())
		}
	}
	if len(errMsg) != 0 {
		return errors.New(strings.Join(errMsg, ";"))
	}
	return nil
}
