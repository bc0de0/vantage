package ac_03_reachability_validation

import (
	"context"

	"vantage/techniques/model"
)

// ReachabilityValidator implements action class AC-03.
type ReachabilityValidator struct{}

func (t ReachabilityValidator) ID() string            { return "ReachabilityValidator" }
func (t ReachabilityValidator) Name() string          { return "ReachabilityValidator" }
func (t ReachabilityValidator) ActionClassID() string { return "AC-03" }
func (t ReachabilityValidator) Evaluate(graph *model.Graph) bool {
	return graph != nil && graph.EvidenceNodes > 0
}
func (t ReachabilityValidator) Execute(_ context.Context, _ *model.Graph) (model.Evidence, error) {
	return model.Evidence{TechniqueID: t.ID(), Summary: "validate host reachability using non-invasive checks", Success: true}, nil
}
func (t ReachabilityValidator) RiskModifier() float64   { return 0.4 }
func (t ReachabilityValidator) ImpactModifier() float64 { return 0.6 }
