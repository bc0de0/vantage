package model

import "context"

// Graph is a read-only snapshot of reasoning graph state used during technique relevance checks.
type Graph struct {
	EvidenceNodes   int
	HypothesisNodes int
	TechniqueNodes  int
	HasSupportsEdge bool
	HasEnablesEdge  bool
}

// Evidence captures normalized output produced by a technique execution.
type Evidence struct {
	TechniqueID string
	Summary     string
	Success     bool
}

// Technique defines the strict contract for MECE technique implementations.
type Technique interface {
	ID() string
	Name() string
	ActionClassID() string
	Evaluate(graph *Graph) bool
	Execute(ctx context.Context, graph *Graph) (Evidence, error)
	RiskModifier() float64
	ImpactModifier() float64
}
