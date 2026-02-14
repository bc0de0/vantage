package memory

import "sync"

// PatternStore holds cross-campaign global archetypes.
type PatternStore struct {
	mu       sync.RWMutex
	patterns map[string]float64
}

func NewPatternStore() *PatternStore {
	return &PatternStore{patterns: map[string]float64{}}
}

func (s *PatternStore) Upsert(name string, confidence float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.patterns[name] = confidence
}
