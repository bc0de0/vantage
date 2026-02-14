package ac_13_data_exposure

import (
	"context"

	"vantage/techniques/model"
)

// DataExposureVerifier implements action class AC-13.
type DataExposureVerifier struct{}

func (t DataExposureVerifier) ID() string            { return "DataExposureVerifier" }
func (t DataExposureVerifier) Name() string          { return "DataExposureVerifier" }
func (t DataExposureVerifier) ActionClassID() string { return "AC-13" }
func (t DataExposureVerifier) Evaluate(graph *model.Graph) bool {
	return graph != nil && graph.HypothesisNodes > 0
}
func (t DataExposureVerifier) Execute(_ context.Context, _ *model.Graph) (model.Evidence, error) {
	return model.Evidence{TechniqueID: t.ID(), Summary: "verify accessible sensitive data exposure", Success: true}, nil
}
func (t DataExposureVerifier) RiskModifier() float64   { return 0.7 }
func (t DataExposureVerifier) ImpactModifier() float64 { return 0.95 }
