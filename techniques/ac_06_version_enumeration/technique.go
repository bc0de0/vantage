package ac_06_version_enumeration

import (
	"context"

	"vantage/techniques/model"
)

// VersionEnumerator implements action class AC-06.
type VersionEnumerator struct{}

func (t VersionEnumerator) ID() string            { return "VersionEnumerator" }
func (t VersionEnumerator) Name() string          { return "VersionEnumerator" }
func (t VersionEnumerator) ActionClassID() string { return "AC-06" }
func (t VersionEnumerator) Evaluate(graph *model.Graph) bool {
	return graph != nil && graph.EvidenceNodes > 0
}
func (t VersionEnumerator) Execute(_ context.Context, _ *model.Graph) (model.Evidence, error) {
	return model.Evidence{TechniqueID: t.ID(), Summary: "enumerate service versions and capabilities", Success: true}, nil
}
func (t VersionEnumerator) RiskModifier() float64   { return 0.6 }
func (t VersionEnumerator) ImpactModifier() float64 { return 0.8 }
