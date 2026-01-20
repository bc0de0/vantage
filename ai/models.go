package ai

// -----------------------------------------------------------------------------
// AI ADVISORY MODELS — NON-AUTHORITATIVE BY DESIGN
//
// This file defines the ONLY data structures that may cross the boundary
// between VANTAGE core logic and the AI advisory layer.
//
// These models are intentionally:
// - Narrow
// - Declarative
// - Tool-agnostic
// - Execution-free
//
// IMPORTANT DOCTRINE:
// - AI NEVER receives raw execution data
// - AI NEVER receives credentials, commands, or steps
// - AI output is advisory only and MUST be ignorable
// -----------------------------------------------------------------------------

// AdvisoryInput represents the full, bounded context
// that the AI is allowed to reason about.
//
// This structure is constructed by deterministic code.
// AI must not infer or guess missing fields.
type AdvisoryInput struct {

	// Intent captures the declared purpose of the campaign.
	Intent struct {
		// Objective is a human-readable description of intent.
		Objective string `json:"objective"`

		// AllowedDomains restricts which intent domains are valid.
		// Example: discovery, enumeration
		AllowedDomains []string `json:"allowed_domains"`
	} `json:"intent"`

	// ROE captures the effective Rules of Engagement.
	ROE struct {
		// AllowedCategories are ROE-permitted action categories.
		AllowedCategories []string `json:"allowed_categories"`

		// ForbiddenCategories are explicitly disallowed.
		ForbiddenCategories []string `json:"forbidden_categories"`
	} `json:"roe"`

	// Exposure represents remaining risk budget.
	Exposure struct {
		// RemainingBudget is a coarse qualitative label.
		// Example: low, medium, high
		RemainingBudget string `json:"remaining_budget"`
	} `json:"exposure"`

	// TargetContext provides high-level, non-sensitive context.
	TargetContext struct {
		// HighLevelType describes the target category.
		// Example: network_service, web_application
		HighLevelType string `json:"high_level_type"`

		// KnownAccess indicates whether access already exists.
		KnownAccess bool `json:"known_access"`
	} `json:"target_context"`

	// CanonicalActionClasses is the complete list of
	// action class IDs AI is allowed to reference.
	CanonicalActionClasses []string `json:"canonical_action_classes"`
}

// AdvisoryOutput is the AI’s advisory-only classification result.
//
// This output MUST:
// - be explainable
// - include exclusions
// - include a mandatory disclaimer
//
// This output MUST NOT:
// - be treated as authoritative
// - drive execution automatically
type AdvisoryOutput struct {

	// Suggested lists action classes deemed relevant.
	Suggested []SuggestedClass `json:"suggested_action_classes"`

	// Excluded lists action classes explicitly ruled out.
	Excluded []ExcludedClass `json:"excluded_action_classes"`

	// Novelty flags indicate patterns that do not map cleanly.
	Novelty []NoveltyFlag `json:"novelty_flags"`

	// Disclaimer MUST be present verbatim.
	Disclaimer string `json:"disclaimer"`
}

// SuggestedClass represents an AI-suggested action class.
type SuggestedClass struct {
	ID         string  `json:"id"`
	Confidence float64 `json:"confidence"`
	Rationale  string  `json:"rationale"`
}

// ExcludedClass represents an explicitly excluded action class.
type ExcludedClass struct {
	ID     string `json:"id"`
	Reason string `json:"reason"`
}

// NoveltyFlag indicates potential gaps in the canonical lattice.
type NoveltyFlag struct {
	Description    string   `json:"description"`
	ClosestMatches []string `json:"closest_matches"`
}
