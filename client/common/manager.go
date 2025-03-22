package common

import "sync"

type manager struct {
	data map[string]map[string]bool
	lock *sync.Mutex
}

var Manager = manager{
	data: make(map[string]map[string]bool),
	lock: &sync.Mutex{},
}

func (m *manager) Store(scope string, key string) {
	data := m.data[scope]
	if data == nil {
		data = make(map[string]bool)
	}
	data[key] = true
	m.data[scope] = data
}

func (m *manager) Remove(scope string, key string) {
	data := m.data[scope]
	if data == nil {
		return
	}
	delete(data, key)
	m.data[scope] = data
}
