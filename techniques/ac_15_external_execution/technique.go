package ac_15_external_execution

import (
	"context"

	"vantage/techniques/model"
)

// ExternalExecutionCoordinator implements action class AC-15.
type ExternalExecutionCoordinator struct{}

func (t ExternalExecutionCoordinator) ID() string            { return "ExternalExecutionCoordinator" }
func (t ExternalExecutionCoordinator) Name() string          { return "ExternalExecutionCoordinator" }
func (t ExternalExecutionCoordinator) ActionClassID() string { return "AC-15" }
func (t ExternalExecutionCoordinator) Evaluate(graph *model.Graph) bool {
	return graph != nil && graph.HasEnablesEdge
}
func (t ExternalExecutionCoordinator) Execute(_ context.Context, _ *model.Graph) (model.Evidence, error) {
	return model.Evidence{TechniqueID: t.ID(), Summary: "coordinate actions requiring external execution dependencies", Success: true}, nil
}
func (t ExternalExecutionCoordinator) RiskModifier() float64   { return 0.75 }
func (t ExternalExecutionCoordinator) ImpactModifier() float64 { return 0.8 }
