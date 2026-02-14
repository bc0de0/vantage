package ac_14_impact_feasibility

import (
	"testing"

	"vantage/techniques/model"
)

func TestImpactFeasibilityAssessorActionClassID(t *testing.T) {
	tech := ImpactFeasibilityAssessor{}
	if tech.ActionClassID() != "AC-14" {
		t.Fatalf("unexpected action class: %s", tech.ActionClassID())
	}
}

func TestImpactFeasibilityAssessorEvaluate(t *testing.T) {
	tech := ImpactFeasibilityAssessor{}
	g := &model.Graph{}
	if tech.Evaluate(g) {
		t.Fatalf("expected irrelevance on empty snapshot")
	}
	g = &model.Graph{EvidenceNodes: 1, HypothesisNodes: 1, TechniqueNodes: 1, HasSupportsEdge: true, HasEnablesEdge: true}
	if !tech.Evaluate(g) {
		t.Fatalf("expected relevance for populated snapshot")
	}
}
