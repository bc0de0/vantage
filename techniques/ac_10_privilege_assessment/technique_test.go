package ac_10_privilege_assessment

import (
	"testing"

	"vantage/techniques/model"
)

func TestPrivilegeAssessorActionClassID(t *testing.T) {
	tech := PrivilegeAssessor{}
	if tech.ActionClassID() != "AC-10" {
		t.Fatalf("unexpected action class: %s", tech.ActionClassID())
	}
}

func TestPrivilegeAssessorEvaluate(t *testing.T) {
	tech := PrivilegeAssessor{}
	g := &model.Graph{}
	if tech.Evaluate(g) {
		t.Fatalf("expected irrelevance on empty snapshot")
	}
	g = &model.Graph{EvidenceNodes: 1, HypothesisNodes: 1, TechniqueNodes: 1, HasSupportsEdge: true, HasEnablesEdge: true}
	if !tech.Evaluate(g) {
		t.Fatalf("expected relevance for populated snapshot")
	}
}
