package ac_09_access_establishment

import (
	"context"

	"vantage/techniques/model"
)

// AccessEstablisher implements action class AC-09.
type AccessEstablisher struct{}

func (t AccessEstablisher) ID() string            { return "AccessEstablisher" }
func (t AccessEstablisher) Name() string          { return "AccessEstablisher" }
func (t AccessEstablisher) ActionClassID() string { return "AC-09" }
func (t AccessEstablisher) Evaluate(graph *model.Graph) bool {
	return graph != nil && graph.TechniqueNodes > 0
}
func (t AccessEstablisher) Execute(_ context.Context, _ *model.Graph) (model.Evidence, error) {
	return model.Evidence{TechniqueID: t.ID(), Summary: "establish authenticated access with validated material", Success: true}, nil
}
func (t AccessEstablisher) RiskModifier() float64   { return 0.9 }
func (t AccessEstablisher) ImpactModifier() float64 { return 0.9 }
