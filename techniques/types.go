package techniques

// ============================================================================
// RESOLUTION TYPES — STRUCTURED, AUDITABLE DECISIONS
// ============================================================================

// ResolveInput is a read-only snapshot of constraints.
//
// ALL enforcement has already occurred.
// Techniques MUST NOT reinterpret policy.
type ResolveInput struct {

	// TechniqueID is redundant but included for safety.
	TechniqueID string

	// AllowedIntentDomains restrict intent scope.
	AllowedIntentDomains []string

	// AllowedROECategories restrict action categories.
	AllowedROECategories []string

	// ExposureBudget is a qualitative label.
	// Example: low, medium, high
	ExposureBudget string

	// CampaignState is informational only.
	// Example: running, halted
	CampaignState string
}

// Resolution is the complete, auditable output of a technique decision.
//
// Silence is NOT allowed.
// A resolution must always explain itself.
type Resolution struct {

	// TechniqueID binds this decision to a technique.
	TechniqueID string

	// AllowedActionClasses are admissible for consideration.
	AllowedActionClasses []string

	// ExcludedActionClasses document explicit denials.
	ExcludedActionClasses []ExcludedActionClass

	// Rationale is a human-readable justification.
	Rationale string
}

// ExcludedActionClass documents why an action class is disallowed.
type ExcludedActionClass struct {
	ActionClassID string
	Reason        string
}
