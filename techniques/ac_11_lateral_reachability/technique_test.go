package ac_11_lateral_reachability

import (
	"testing"

	"vantage/techniques/model"
)

func TestLateralReachabilityAnalyzerActionClassID(t *testing.T) {
	tech := LateralReachabilityAnalyzer{}
	if tech.ActionClassID() != "AC-11" {
		t.Fatalf("unexpected action class: %s", tech.ActionClassID())
	}
}

func TestLateralReachabilityAnalyzerEvaluate(t *testing.T) {
	tech := LateralReachabilityAnalyzer{}
	g := &model.Graph{}
	if tech.Evaluate(g) {
		t.Fatalf("expected irrelevance on empty snapshot")
	}
	g = &model.Graph{EvidenceNodes: 1, HypothesisNodes: 1, TechniqueNodes: 1, HasSupportsEdge: true, HasEnablesEdge: true}
	if !tech.Evaluate(g) {
		t.Fatalf("expected relevance for populated snapshot")
	}
}
