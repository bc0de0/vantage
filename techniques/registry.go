package techniques

import (
	"fmt"
	"sync"
)

// ============================================================================
// TECHNIQUE REGISTRY — SECURITY CRITICAL COMPONENT
//
// This file defines the ONLY authoritative registry of techniques in VANTAGE.
//
// GUARANTEES:
// - Closed world (no dynamic discovery)
// - Immutable after init()
// - Fail-fast on programmer error
// - Deterministic lookup
//
// If a technique is not registered here at process startup,
// it does not exist — structurally or legally.
// ============================================================================

var (
	// registry maps canonical Technique IDs to implementations.
	// It is intentionally unexported.
	registry = make(map[string]Technique)

	// registryMu protects registry during init-time registration.
	// Runtime writes are forbidden by doctrine.
	registryMu sync.RWMutex
)

// Register registers a Technique implementation.
//
// HARD RULES (ENFORCED):
// - MUST be called only from init()
// - MUST NOT be called conditionally
// - MUST NOT be called at runtime
//
// Any violation indicates a programmer error and causes panic.
func Register(t Technique) {
	if t == nil {
		panic("techniques: attempted to register nil Technique")
	}

	id := t.ID()
	if id == "" {
		panic("techniques: attempted to register Technique with empty ID")
	}

	registryMu.Lock()
	defer registryMu.Unlock()

	if _, exists := registry[id]; exists {
		panic("techniques: duplicate technique registration: " + id)
	}

	registry[id] = t
}

// Get retrieves a registered Technique by ID.
//
// BEHAVIOR:
// - Returns error if technique does not exist
// - Never panics for missing techniques
//
// Missing techniques are an operator/configuration error,
// not a programmer error.
func Get(id string) (Technique, error) {
	if id == "" {
		return nil, fmt.Errorf("techniques: empty technique ID")
	}

	registryMu.RLock()
	defer registryMu.RUnlock()

	t, ok := registry[id]
	if !ok {
		return nil, fmt.Errorf("techniques: technique not registered: %s", id)
	}

	return t, nil
}

// List returns all registered technique IDs.
//
// Intended for:
// - CLI validation
// - Intent verification
// - Audit tooling
func List() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()

	out := make([]string, 0, len(registry))
	for id := range registry {
		out = append(out, id)
	}
	return out
}
