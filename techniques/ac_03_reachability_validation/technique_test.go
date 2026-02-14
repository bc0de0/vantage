package ac_03_reachability_validation

import (
	"testing"

	"vantage/techniques/model"
)

func TestReachabilityValidatorActionClassID(t *testing.T) {
	tech := ReachabilityValidator{}
	if tech.ActionClassID() != "AC-03" {
		t.Fatalf("unexpected action class: %s", tech.ActionClassID())
	}
}

func TestReachabilityValidatorEvaluate(t *testing.T) {
	tech := ReachabilityValidator{}
	g := &model.Graph{}
	if tech.Evaluate(g) {
		t.Fatalf("expected irrelevance on empty snapshot")
	}
	g = &model.Graph{EvidenceNodes: 1, HypothesisNodes: 1, TechniqueNodes: 1, HasSupportsEdge: true, HasEnablesEdge: true}
	if !tech.Evaluate(g) {
		t.Fatalf("expected relevance for populated snapshot")
	}
}
