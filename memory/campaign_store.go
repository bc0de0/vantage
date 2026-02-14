package memory

import "sync"

// CampaignStore persists reusable campaign patterns.
type CampaignStore struct {
	mu       sync.RWMutex
	patterns map[string][]string
}

func NewCampaignStore() *CampaignStore {
	return &CampaignStore{patterns: map[string][]string{}}
}

func (s *CampaignStore) Record(campaignID string, findings []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.patterns[campaignID] = append([]string(nil), findings...)
}
