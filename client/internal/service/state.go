package service

import (
	"sync"
)

var State = RunningState{
	mu:   &sync.RWMutex{},
	data: make(map[string]bool),
}

type RunningState struct {
	mu   *sync.RWMutex
	data map[string]bool
}

func (state *RunningState) Set(key string, running bool) {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.data[key] = running
}

func (state *RunningState) Get(key string) bool {
	state.mu.RLock()
	defer state.mu.RUnlock()
	return state.data[key]
}
