package database

import (
	"sync"
)

type AccessCounter struct {
	counts map[string]int
	mu     sync.Mutex
}

func NewAccessCounter() *AccessCounter {
	return &AccessCounter{
		counts: make(map[string]int),
	}
}

func (ac *AccessCounter) Increment(key string) int {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.counts[key]++
	return ac.counts[key]
}

func (ac *AccessCounter) Get(key string) int {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	return ac.counts[key]
}

func (ac *AccessCounter) Reset(key string) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	delete(ac.counts, key)
}
