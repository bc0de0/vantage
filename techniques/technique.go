package techniques

// ============================================================================
// TECHNIQUE — DECISION CONTRACT (SECURITY CRITICAL)
//
// A Technique in VANTAGE:
// - DOES NOT execute actions
// - DOES NOT touch the network
// - DOES NOT mutate state
// - DOES NOT call AI
//
// A Technique ONLY decides which Action Classes are admissible,
// given constraints already enforced by the Engine.
// ============================================================================

// Technique is the canonical decision interface.
type Technique interface {

	// ID returns the stable, canonical technique identifier.
	//
	// FORMAT:
	// - Uppercase
	// - MITRE-style (e.g., "T1595", "T1078.003")
	//
	// This value is part of the security boundary.
	ID() string

	// Description returns a human-readable summary.
	//
	// This is used for:
	// - Audit output
	// - Reports
	// - Non-technical review
	//
	// MUST be non-operational.
	Description() string

	// Resolve determines admissible Action Classes.
	//
	// This function MUST be:
	// - Pure
	// - Deterministic
	// - Side-effect free
	//
	// Errors indicate structural failure,
	// NOT policy denial (policy is enforced by the engine).
	Resolve(input ResolveInput) (*Resolution, error)
}
