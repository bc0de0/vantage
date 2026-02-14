package ac_13_data_exposure

import (
	"testing"

	"vantage/techniques/model"
)

func TestDataExposureVerifierActionClassID(t *testing.T) {
	tech := DataExposureVerifier{}
	if tech.ActionClassID() != "AC-13" {
		t.Fatalf("unexpected action class: %s", tech.ActionClassID())
	}
}

func TestDataExposureVerifierEvaluate(t *testing.T) {
	tech := DataExposureVerifier{}
	g := &model.Graph{}
	if tech.Evaluate(g) {
		t.Fatalf("expected irrelevance on empty snapshot")
	}
	g = &model.Graph{EvidenceNodes: 1, HypothesisNodes: 1, TechniqueNodes: 1, HasSupportsEdge: true, HasEnablesEdge: true}
	if !tech.Evaluate(g) {
		t.Fatalf("expected relevance for populated snapshot")
	}
}
