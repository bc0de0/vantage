package ac_01_passive_observation

import (
	"testing"

	"vantage/techniques/model"
)

func TestPassiveDNSCollectionActionClassID(t *testing.T) {
	tech := PassiveDNSCollection{}
	if tech.ActionClassID() != "AC-01" {
		t.Fatalf("unexpected action class: %s", tech.ActionClassID())
	}
}

func TestPassiveDNSCollectionEvaluate(t *testing.T) {
	tech := PassiveDNSCollection{}
	g := &model.Graph{}
	if !tech.Evaluate(g) {
		t.Fatalf("expected relevance with empty evidence")
	}
	g = &model.Graph{EvidenceNodes: 1, HypothesisNodes: 1, TechniqueNodes: 1, HasSupportsEdge: true, HasEnablesEdge: true}
	if tech.Evaluate(g) {
		t.Fatalf("expected irrelevance after evidence appears")
	}
}
