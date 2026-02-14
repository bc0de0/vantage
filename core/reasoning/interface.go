package reasoning

import (
	"vantage/core/evidence"
	"vantage/core/state"
)

// TechniqueEffect describes expected outcomes for a technique.
type TechniqueEffect struct {
	TechniqueID   string
	ActionClassID string
	Impact        float64
	Risk          float64
	Stealth       float64
	Produces      []string
}

// TechniqueEffectRegistry stores technique effects used during planning.
type TechniqueEffectRegistry interface {
	RegisterTechniqueEffect(effect TechniqueEffect)
	EffectForTechnique(techniqueID string) (TechniqueEffect, bool)
	KnownTechniques() []string
}

// EvidenceEvent is a normalized event emitted by the executor.
type EvidenceEvent struct {
	TechniqueID string
	Target      string
	Success     bool
	Output      string
	Artifact    *evidence.Artifact
}

// EvidenceIngestor accepts evidence events for reasoning updates.
type EvidenceIngestor interface {
	IngestEvidence(event EvidenceEvent) error
}

// PlannerQuery is the planner query input.
type PlannerQuery struct {
	Target             string
	AllowedTechniques  []string
	CurrentTechniqueID string
	TopN               int
}

// RankedAction is a scored action candidate returned by the planner.
type RankedAction struct {
	TechniqueID   string
	ActionClassID string
	Target        string
	Score         float64
	Impact        float64
	Risk          float64
	Stealth       float64
	Reason        string
}

// RankedActionPlanner returns ranked next actions.
type RankedActionPlanner interface {
	RankedActions(query PlannerQuery) []RankedAction
}

// HypothesisExpander allows optional advisory hypothesis generation.
type HypothesisExpander interface {
	Expand(graph *Graph, state *state.State) ([]Hypothesis, error)
}
