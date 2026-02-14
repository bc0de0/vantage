package ac_02_active_surface_discovery

import (
	"context"

	"vantage/techniques/model"
)

// SurfaceProbe implements action class AC-02.
type SurfaceProbe struct{}

func (t SurfaceProbe) ID() string            { return "SurfaceProbe" }
func (t SurfaceProbe) Name() string          { return "SurfaceProbe" }
func (t SurfaceProbe) ActionClassID() string { return "AC-02" }
func (t SurfaceProbe) Evaluate(graph *model.Graph) bool {
	return graph != nil && graph.EvidenceNodes > 0
}
func (t SurfaceProbe) Execute(_ context.Context, _ *model.Graph) (model.Evidence, error) {
	return model.Evidence{TechniqueID: t.ID(), Summary: "probe target surface for reachable hosts and endpoints", Success: true}, nil
}
func (t SurfaceProbe) RiskModifier() float64   { return 0.5 }
func (t SurfaceProbe) ImpactModifier() float64 { return 0.6 }
