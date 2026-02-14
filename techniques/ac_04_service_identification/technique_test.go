package ac_04_service_identification

import (
	"testing"

	"vantage/techniques/model"
)

func TestServiceIdentifierActionClassID(t *testing.T) {
	tech := ServiceIdentifier{}
	if tech.ActionClassID() != "AC-04" {
		t.Fatalf("unexpected action class: %s", tech.ActionClassID())
	}
}

func TestServiceIdentifierEvaluate(t *testing.T) {
	tech := ServiceIdentifier{}
	g := &model.Graph{}
	if tech.Evaluate(g) {
		t.Fatalf("expected irrelevance on empty snapshot")
	}
	g = &model.Graph{EvidenceNodes: 1, HypothesisNodes: 1, TechniqueNodes: 1, HasSupportsEdge: true, HasEnablesEdge: true}
	if !tech.Evaluate(g) {
		t.Fatalf("expected relevance for populated snapshot")
	}
}
