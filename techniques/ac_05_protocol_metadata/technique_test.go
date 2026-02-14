package ac_05_protocol_metadata

import (
	"testing"

	"vantage/techniques/model"
)

func TestProtocolMetadataInspectorActionClassID(t *testing.T) {
	tech := ProtocolMetadataInspector{}
	if tech.ActionClassID() != "AC-05" {
		t.Fatalf("unexpected action class: %s", tech.ActionClassID())
	}
}

func TestProtocolMetadataInspectorEvaluate(t *testing.T) {
	tech := ProtocolMetadataInspector{}
	g := &model.Graph{}
	if tech.Evaluate(g) {
		t.Fatalf("expected irrelevance on empty snapshot")
	}
	g = &model.Graph{EvidenceNodes: 1, HypothesisNodes: 1, TechniqueNodes: 1, HasSupportsEdge: true, HasEnablesEdge: true}
	if !tech.Evaluate(g) {
		t.Fatalf("expected relevance for populated snapshot")
	}
}
