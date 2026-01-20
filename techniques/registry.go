package techniques

import (
	"fmt"
	"sync"
)

// -----------------------------------------------------------------------------
// Technique Registry — SECURITY CRITICAL COMPONENT
//
// This file defines the ONLY authoritative registry of executable techniques
// within VANTAGE.
//
// DESIGN PRINCIPLES:
// 1. CLOSED WORLD ASSUMPTION
//    Only techniques explicitly registered at process startup exist.
//    If it is not registered here, it does not exist.
//
// 2. IMMUTABILITY AFTER INIT
//    Registration is allowed ONLY during program initialization.
//    Runtime mutation is forbidden.
//
// 3. FAIL FAST ON VIOLATION
//    Duplicate registrations or invalid technique definitions are
//    programmer errors and MUST panic.
//
// 4. NO STRINGLY-TYPED EXECUTION
//    Techniques are resolved to concrete implementations before execution.
// -----------------------------------------------------------------------------

// registry holds all known techniques keyed by technique ID (e.g., "T1595").
// It is intentionally unexported to prevent external mutation.
//
// Access is guarded by a mutex to protect against:
// - accidental concurrent registration
// - future refactors that introduce parallel init paths
var (
	registry   = make(map[string]Technique)
	registryMu sync.RWMutex
)

// Register registers a Technique implementation with the global registry.
//
// STRICT RULES:
// - Must ONLY be called from init() functions inside technique packages
// - Must NEVER be called at runtime
// - Must NEVER be called conditionally
//
// Violations of these rules indicate a developer error and will PANIC.
//
// This is intentional: silent failure here would compromise ROE enforcement.
func Register(t Technique) {
	if t == nil {
		panic("attempted to register nil Technique")
	}

	id := t.ID()

	// Technique IDs are part of the security boundary.
	// Empty or malformed IDs are not acceptable.
	if id == "" {
		panic("attempted to register technique with empty ID")
	}

	registryMu.Lock()
	defer registryMu.Unlock()

	// Duplicate technique IDs are a hard failure.
	// Allowing overrides would enable technique shadowing.
	if _, exists := registry[id]; exists {
		panic("duplicate technique registration detected: " + id)
	}

	registry[id] = t
}

// Get retrieves a registered Technique by ID.
//
// This function is used by the executor and is the ONLY supported
// lookup mechanism.
//
// BEHAVIOR:
// - Returns a concrete Technique if registered
// - Returns an explicit error if the technique does not exist
//
// NOTE:
//   - This function does NOT panic on missing techniques.
//   - Missing techniques are an operator/configuration error,
//     not a programmer error.
func Get(id string) (Technique, error) {
	if id == "" {
		return nil, fmt.Errorf("invalid technique ID: empty string")
	}

	registryMu.RLock()
	defer registryMu.RUnlock()

	t, ok := registry[id]
	if !ok {
		return nil, fmt.Errorf("technique not registered: %s", id)
	}

	return t, nil
}

// List returns the set of registered technique IDs.
//
// This is intentionally read-only and is provided for:
// - CLI validation
// - intent contract verification
// - audit and inspection commands
//
// Callers MUST NOT assume ordering.
func List() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()

	ids := make([]string, 0, len(registry))
	for id := range registry {
		ids = append(ids, id)
	}

	return ids
}
