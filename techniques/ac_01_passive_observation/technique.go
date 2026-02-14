package ac_01_passive_observation

import (
	"context"

	"vantage/techniques/model"
)

// PassiveDNSCollection implements action class AC-01.
type PassiveDNSCollection struct{}

func (t PassiveDNSCollection) ID() string            { return "PassiveDNSCollection" }
func (t PassiveDNSCollection) Name() string          { return "PassiveDNSCollection" }
func (t PassiveDNSCollection) ActionClassID() string { return "AC-01" }
func (t PassiveDNSCollection) Evaluate(graph *model.Graph) bool {
	return graph != nil && graph.EvidenceNodes == 0
}
func (t PassiveDNSCollection) Execute(_ context.Context, _ *model.Graph) (model.Evidence, error) {
	return model.Evidence{TechniqueID: t.ID(), Summary: "collect passive DNS and certificate transparency observations", Success: true}, nil
}
func (t PassiveDNSCollection) RiskModifier() float64   { return 0.2 }
func (t PassiveDNSCollection) ImpactModifier() float64 { return 0.4 }
