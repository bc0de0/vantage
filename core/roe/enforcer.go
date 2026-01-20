package roe

import (
	"errors"
	"fmt"
	"time"

	"vantage/core/intent"
)

// -----------------------------------------------------------------------------
// Rules of Engagement (ROE) Enforcer — AUTHORITATIVE POLICY GATE
//
// ROE is enforced as INTERSECTION, not override.
// Execution is permitted only when BOTH:
//   - Static ROE allows it
//   - Declared intent allows it
//
// Absence of permission = denial.
// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------
// STATIC ROE (v0.x — COMPILED POLICY)
// -----------------------------------------------------------------------------

var allowedTechniques = map[string]struct{}{
	"T1595": {}, // Active Scanning
}

// -----------------------------------------------------------------------------
// Enforce — FINAL AUTHORITY
// -----------------------------------------------------------------------------

// Enforce validates whether a technique execution is permitted
// under BOTH ROE and declared intent.
//
// This function is:
// - deterministic
// - side-effect free
// - fail-closed
func Enforce(contract *intent.Contract, techniqueID string, target string) error {

	// -----------------------------
	// 1. Defensive Validation
	// -----------------------------

	if contract == nil {
		return errors.New("ROE violation: nil intent contract")
	}

	if techniqueID == "" {
		return errors.New("ROE violation: empty technique ID")
	}

	if target == "" {
		return errors.New("ROE violation: empty target")
	}

	// -----------------------------
	// 2. Static ROE Enforcement
	// -----------------------------

	if _, ok := allowedTechniques[techniqueID]; !ok {
		return fmt.Errorf(
			"ROE violation: technique %s not permitted by policy",
			techniqueID,
		)
	}

	// -----------------------------
	// 3. Intent-Declared Technique Scope
	// -----------------------------

	allowed := false
	for _, t := range contract.AllowedTechniques {
		if t == techniqueID {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf(
			"ROE violation: technique %s not declared in intent",
			techniqueID,
		)
	}

	// -----------------------------
	// 4. Intent-Declared Target Scope
	// -----------------------------

	targetAllowed := false
	for _, t := range contract.Targets {
		if t == target {
			targetAllowed = true
			break
		}
	}

	if !targetAllowed {
		return fmt.Errorf(
			"ROE violation: target %s not declared in intent",
			target,
		)
	}

	// -----------------------------
	// 5. Temporal Enforcement
	// -----------------------------

	now := time.Now().UTC()

	if now.Before(contract.NotBefore) || now.After(contract.NotAfter) {
		return fmt.Errorf(
			"ROE violation: execution outside intent window (UTC %s → %s)",
			contract.NotBefore.Format(time.RFC3339),
			contract.NotAfter.Format(time.RFC3339),
		)
	}

	// -----------------------------
	// PASS — EXECUTION PERMITTED
	// -----------------------------
	return nil
}
