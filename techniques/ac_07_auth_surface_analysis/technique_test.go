package ac_07_auth_surface_analysis

import (
	"testing"

	"vantage/techniques/model"
)

func TestAuthSurfaceAnalyzerActionClassID(t *testing.T) {
	tech := AuthSurfaceAnalyzer{}
	if tech.ActionClassID() != "AC-07" {
		t.Fatalf("unexpected action class: %s", tech.ActionClassID())
	}
}

func TestAuthSurfaceAnalyzerEvaluate(t *testing.T) {
	tech := AuthSurfaceAnalyzer{}
	g := &model.Graph{}
	if tech.Evaluate(g) {
		t.Fatalf("expected irrelevance on empty snapshot")
	}
	g = &model.Graph{EvidenceNodes: 1, HypothesisNodes: 1, TechniqueNodes: 1, HasSupportsEdge: true, HasEnablesEdge: true}
	if !tech.Evaluate(g) {
		t.Fatalf("expected relevance for populated snapshot")
	}
}
