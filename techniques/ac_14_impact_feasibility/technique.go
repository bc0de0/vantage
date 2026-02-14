package ac_14_impact_feasibility

import (
	"context"

	"vantage/techniques/model"
)

// ImpactFeasibilityAssessor implements action class AC-14.
type ImpactFeasibilityAssessor struct{}

func (t ImpactFeasibilityAssessor) ID() string            { return "ImpactFeasibilityAssessor" }
func (t ImpactFeasibilityAssessor) Name() string          { return "ImpactFeasibilityAssessor" }
func (t ImpactFeasibilityAssessor) ActionClassID() string { return "AC-14" }
func (t ImpactFeasibilityAssessor) Evaluate(graph *model.Graph) bool {
	return graph != nil && graph.HypothesisNodes > 0
}
func (t ImpactFeasibilityAssessor) Execute(_ context.Context, _ *model.Graph) (model.Evidence, error) {
	return model.Evidence{TechniqueID: t.ID(), Summary: "assess whether operational impact is feasible", Success: true}, nil
}
func (t ImpactFeasibilityAssessor) RiskModifier() float64   { return 0.9 }
func (t ImpactFeasibilityAssessor) ImpactModifier() float64 { return 0.95 }
