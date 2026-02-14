package ac_05_protocol_metadata

import (
	"context"

	"vantage/techniques/model"
)

// ProtocolMetadataInspector implements action class AC-05.
type ProtocolMetadataInspector struct{}

func (t ProtocolMetadataInspector) ID() string            { return "ProtocolMetadataInspector" }
func (t ProtocolMetadataInspector) Name() string          { return "ProtocolMetadataInspector" }
func (t ProtocolMetadataInspector) ActionClassID() string { return "AC-05" }
func (t ProtocolMetadataInspector) Evaluate(graph *model.Graph) bool {
	return graph != nil && graph.EvidenceNodes > 0
}
func (t ProtocolMetadataInspector) Execute(_ context.Context, _ *model.Graph) (model.Evidence, error) {
	return model.Evidence{TechniqueID: t.ID(), Summary: "inspect protocol banners and metadata", Success: true}, nil
}
func (t ProtocolMetadataInspector) RiskModifier() float64   { return 0.4 }
func (t ProtocolMetadataInspector) ImpactModifier() float64 { return 0.7 }
