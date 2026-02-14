package ac_04_service_identification

import (
	"context"

	"vantage/techniques/model"
)

// ServiceIdentifier implements action class AC-04.
type ServiceIdentifier struct{}

func (t ServiceIdentifier) ID() string            { return "ServiceIdentifier" }
func (t ServiceIdentifier) Name() string          { return "ServiceIdentifier" }
func (t ServiceIdentifier) ActionClassID() string { return "AC-04" }
func (t ServiceIdentifier) Evaluate(graph *model.Graph) bool {
	return graph != nil && graph.HasSupportsEdge
}
func (t ServiceIdentifier) Execute(_ context.Context, _ *model.Graph) (model.Evidence, error) {
	return model.Evidence{TechniqueID: t.ID(), Summary: "identify exposed network services", Success: true}, nil
}
func (t ServiceIdentifier) RiskModifier() float64   { return 0.5 }
func (t ServiceIdentifier) ImpactModifier() float64 { return 0.7 }
