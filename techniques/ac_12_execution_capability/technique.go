package ac_12_execution_capability

import (
	"context"

	"vantage/techniques/model"
)

// ExecutionCapabilityValidator implements action class AC-12.
type ExecutionCapabilityValidator struct{}

func (t ExecutionCapabilityValidator) ID() string            { return "ExecutionCapabilityValidator" }
func (t ExecutionCapabilityValidator) Name() string          { return "ExecutionCapabilityValidator" }
func (t ExecutionCapabilityValidator) ActionClassID() string { return "AC-12" }
func (t ExecutionCapabilityValidator) Evaluate(graph *model.Graph) bool {
	return graph != nil && graph.HypothesisNodes > 0
}
func (t ExecutionCapabilityValidator) Execute(_ context.Context, _ *model.Graph) (model.Evidence, error) {
	return model.Evidence{TechniqueID: t.ID(), Summary: "validate in-environment execution capability", Success: true}, nil
}
func (t ExecutionCapabilityValidator) RiskModifier() float64   { return 0.85 }
func (t ExecutionCapabilityValidator) ImpactModifier() float64 { return 0.9 }
