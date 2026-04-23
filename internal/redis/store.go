package redis

import (
	"sync"
	"time"
)

type record struct {
	value     []byte
	expiresAt time.Time
}

type list struct {
	first *listEntry
	last  *listEntry
	count int
}

type listEntry struct {
	value    []byte
	previous *listEntry
	next     *listEntry
}

type Store struct {
	mu      sync.RWMutex
	records map[string]*record
	lists   map[string]*list
}

func NewStore() *Store {
	return &Store{
		records: make(map[string]*record),
		lists:   make(map[string]*list),
	}
}

func (r *record) isExpired() bool {
	if r.expiresAt.IsZero() {
		return false
	}
	return time.Now().After(r.expiresAt)
}

func (l *list) append(value []byte) bool {
	newEntry := &listEntry{value: value}
	if l.last == nil {
		l.first = newEntry
	} else {
		newEntry.previous = l.last
		l.last.next = newEntry
	}
	l.last = newEntry
	l.count++
	return true
}

// func (l *list) prepend(value []byte) bool {
// 	newEntry := &listEntry{value: value}
// 	if l.first == nil {
// 		l.last = newEntry
// 	} else {
// 		newEntry.next = l.first
// 		l.first.previous = newEntry
// 	}
// 	l.first = newEntry
// 	l.count++
// 	return true
// }

func (s *Store) Set(key string, data []byte, opts setOptions) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := &record{
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
	if r.isExpired() {
		return nil, false
	}
	return r.value, true
}

func (s *Store) Rpush(key string, values ...[]byte) (int, bool) {
	l, ok := s.lists[key]
	if !ok {
		l = &list{}
		s.lists[key] = l
	}

	for _, value := range values {
		ok := l.append(value)
		if !ok {
			return 0, false
		}
	}
	return l.count, true
}
