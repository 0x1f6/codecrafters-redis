package redis

import (
	"sync"
	"time"
)

type Record struct {
	value     []byte
	expiresAt time.Time
}

type Store struct {
	mu      sync.RWMutex
	records map[string]Record
}

func NewStore() *Store {
	return &Store{
		records: make(map[string]Record),
	}
}

func (r *Record) IsExpired() bool {
	if r.expiresAt.IsZero() {
		return false
	}
	return time.Now().After(r.expiresAt)
}

func (s *Store) Set(key string, data []byte, opts setOptions) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := Record{
		value: data,
	}
	if !opts.expiresAt.IsZero() {
		r.expiresAt = opts.expiresAt
	}

	s.records[key] = r
	return true
}

func (s *Store) Get(key string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	r, ok := s.records[key]
	if !ok {
		return nil, false
	}
	if r.IsExpired() {
		return nil, false
	}
	return r.value, true
}
