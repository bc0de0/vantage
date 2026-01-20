package techniques

import "context"

// -----------------------------------------------------------------------------
// Technique — EXECUTION CONTRACT (SECURITY CRITICAL)
//
// This interface defines the ONLY permissible shape of an executable
// adversary action in VANTAGE.
//
// Any code that performs adversary behavior MUST implement this interface
// and be registered via the technique registry.
//
// This interface is intentionally minimal and intentionally restrictive.
//
// DESIGN DOCTRINE:
// 1. ATOMICITY
//    A Technique represents ONE bounded adversary action.
//    It must not chain other techniques, escalate privileges,
//    or persist access.
//
// 2. DETERMINISM
//    Given the same inputs and environment, behavior must be predictable.
//    Techniques must not depend on hidden global state.
//
// 3. EXPLAINABILITY
//    Every technique must be explainable to:
//    - a CISO
//    - a legal reviewer
//    - a defensive engineer
//
// 4. EXECUTOR OWNERSHIP
//    Techniques do NOT control:
//    - campaign state
//    - ROE enforcement
//    - exposure tracking
//    - evidence persistence
//
//    They perform work. Nothing more.
//
// 5. NO SIDE CHANNELS
//    Techniques must not:
//    - spawn background goroutines
//    - write to disk
//    - open listeners
//    - communicate over the network except for the intended action
//
//    All outputs must flow through the executor.
// -----------------------------------------------------------------------------

// Technique is the canonical interface for all executable adversary actions.
//
// Implementations MUST:
// - Be safe to call once
// - Respect context cancellation
// - Fail fast on error
//
// Implementations MUST NOT:
// - Panic
// - Retry indefinitely
// - Suppress errors
// - Mutate global state
type Technique interface {

	// ID returns the canonical technique identifier.
	//
	// Format:
	// - Uppercase
	// - MITRE-style (e.g., "T1595", "T1078.003")
	//
	// This ID is part of the security boundary and MUST be stable.
	ID() string

	// Description returns a human-readable summary of the technique.
	//
	// This value is used for:
	// - operator visibility
	// - audit logs
	// - reporting
	//
	// It MUST NOT:
	// - contain implementation details
	// - speculate on impact
	// - include dynamic data
	Description() string

	// Execute performs the adversary action against a single target.
	//
	// PARAMETERS:
	// - ctx:
	//   Provided by the executor and MUST be honored.
	//   Cancellation indicates:
	//     - timeout reached
	//     - ROE violation
	//     - campaign halt
	//
	// - target:
	//   A single, executor-validated target identifier.
	//   Interpretation is technique-specific.
	//
	// RETURNS:
	// - error:
	//   A non-nil error indicates the technique did not complete successfully.
	//
	// EXECUTION RULES:
	// - Must block until completion or context cancellation
	// - Must return control promptly when ctx.Done() is closed
	// - Must perform exactly ONE adversary action
	//
	// EVIDENCE:
	// - Techniques do NOT create or persist evidence directly.
	// - All evidence is collected and signed by the executor.
	Execute(ctx context.Context, target string) error
}
