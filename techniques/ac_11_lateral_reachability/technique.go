package ac_11_lateral_reachability

import (
	"context"

	"vantage/techniques/model"
)

// LateralReachabilityAnalyzer implements action class AC-11.
type LateralReachabilityAnalyzer struct{}

func (t LateralReachabilityAnalyzer) ID() string            { return "LateralReachabilityAnalyzer" }
func (t LateralReachabilityAnalyzer) Name() string          { return "LateralReachabilityAnalyzer" }
func (t LateralReachabilityAnalyzer) ActionClassID() string { return "AC-11" }
func (t LateralReachabilityAnalyzer) Evaluate(graph *model.Graph) bool {
	return graph != nil && graph.HasEnablesEdge
}
func (t LateralReachabilityAnalyzer) Execute(_ context.Context, _ *model.Graph) (model.Evidence, error) {
	return model.Evidence{TechniqueID: t.ID(), Summary: "assess lateral movement reachability from established access", Success: true}, nil
}
func (t LateralReachabilityAnalyzer) RiskModifier() float64   { return 0.8 }
func (t LateralReachabilityAnalyzer) ImpactModifier() float64 { return 0.85 }
