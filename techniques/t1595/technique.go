package t1595

import (
	"context"
	"errors"

	"vantage/techniques"
)

// -----------------------------------------------------------------------------
// Technique: T1595 — Active Scanning
//
// STATUS:
// - Stub implementation (v0.x)
// - Execution logic intentionally absent
//
// PURPOSE:
// This file serves as the canonical reference implementation for all
// VANTAGE techniques.
//
// Any future technique MUST follow this structure exactly.
// Deviations require architectural review.
// -----------------------------------------------------------------------------

// Technique1595 implements the techniques.Technique interface.
//
// DESIGN NOTES:
// - Stateless by design
// - No configuration stored on the struct
// - Safe for reuse across executions
//
// Techniques MUST NOT:
// - Carry mutable fields
// - Cache results
// - Store runtime state
//
// All state belongs to the executor.
type Technique1595 struct{}

// ID returns the canonical technique identifier.
//
// This value is:
// - Security-critical
// - Used for ROE enforcement
// - Used for registry lookup
// - Used in evidence and reporting
//
// It MUST:
// - Be stable
// - Match MITRE nomenclature
// - Never change once released
func (t *Technique1595) ID() string {
	return "T1595"
}

// Description returns a human-readable summary of the technique.
//
// This text is used for:
// - Operator awareness
// - Audit output
// - Reports
//
// It MUST:
// - Be factual
// - Be concise
// - Avoid operational detail
// - Avoid speculative impact language
func (t *Technique1595) Description() string {
	return "Active Scanning"
}

// Execute performs the technique against a single target.
//
// IMPORTANT EXECUTION RULES:
//
// 1. This function MUST honor context cancellation.
// 2. This function MUST block until completion or cancellation.
// 3. This function MUST perform exactly ONE adversary action.
// 4. This function MUST NOT:
//   - Spawn goroutines
//   - Retry indefinitely
//   - Write to disk
//   - Persist access
//   - Call other techniques
//
// RETURN VALUE:
// - nil     → Technique executed successfully
// - error   → Technique failed or was aborted
//
// NOTE:
// This stub intentionally returns an error to prevent accidental use.
// Real execution logic will be added only after ROE, exposure, and
// evidence pipelines are finalized.
func (t *Technique1595) Execute(ctx context.Context, target string) error {
	// Defensive guard: context must be non-nil
	if ctx == nil {
		return errors.New("nil context provided to technique")
	}

	// Defensive guard: target must be non-empty
	if target == "" {
		return errors.New("empty target provided to technique")
	}

	// Explicit stub failure.
	// This prevents silent success during early development
	// and forces intentional implementation.
	return errors.New("T1595 execution not implemented")
}

// init registers the technique with the global registry.
//
// REGISTRATION RULES:
// - Registration MUST happen in init()
// - Registration MUST be unconditional
// - Registration MUST panic on failure
//
// This ensures the technique is visible at startup
// and cannot be dynamically injected at runtime.
func init() {
	techniques.Register(&Technique1595{})
}
