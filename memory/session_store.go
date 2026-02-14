package memory

import "sync"

// SessionStore keeps volatile context for a single operation.
type SessionStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewSessionStore() *SessionStore {
	return &SessionStore{data: map[string]string{}}
}

func (s *SessionStore) Put(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *SessionStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	return v, ok
}
