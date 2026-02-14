package reasoning

import "vantage/core/state"

// Hypothesis captures a testable offensive claim with transparent evidence.
type Hypothesis struct {
	Claim             string
	Evidence          []string
	SourceReliability float64
	InferenceDepth    int
	Confidence        float64
	Phase             state.OperationPhase
}

// AttackOption is a ranked plan candidate for the next cycle.
type AttackOption struct {
	Name           string
	Phase          state.OperationPhase
	ImpactScore    float64
	StealthScore   float64
	EffortScore    float64
	DependencyCost float64
	Score          float64
	Rationale      string
}

// CycleInput is the deterministic input to one reasoning iteration.
type CycleInput struct {
	Phase      state.OperationPhase
	Signals    []string
	Hypotheses []Hypothesis
}

// CycleOutput is the complete result produced by the reasoning core.
type CycleOutput struct {
	SelectedPhase state.OperationPhase
	TopOptions    []AttackOption
	Hypotheses    []Hypothesis
}
