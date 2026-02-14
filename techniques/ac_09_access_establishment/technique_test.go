package ac_09_access_establishment

import (
	"testing"

	"vantage/techniques/model"
)

func TestAccessEstablisherActionClassID(t *testing.T) {
	tech := AccessEstablisher{}
	if tech.ActionClassID() != "AC-09" {
		t.Fatalf("unexpected action class: %s", tech.ActionClassID())
	}
}

func TestAccessEstablisherEvaluate(t *testing.T) {
	tech := AccessEstablisher{}
	g := &model.Graph{}
	if tech.Evaluate(g) {
		t.Fatalf("expected irrelevance on empty snapshot")
	}
	g = &model.Graph{EvidenceNodes: 1, HypothesisNodes: 1, TechniqueNodes: 1, HasSupportsEdge: true, HasEnablesEdge: true}
	if !tech.Evaluate(g) {
		t.Fatalf("expected relevance for populated snapshot")
	}
}
