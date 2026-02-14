package ac_06_version_enumeration

import (
	"testing"

	"vantage/techniques/model"
)

func TestVersionEnumeratorActionClassID(t *testing.T) {
	tech := VersionEnumerator{}
	if tech.ActionClassID() != "AC-06" {
		t.Fatalf("unexpected action class: %s", tech.ActionClassID())
	}
}

func TestVersionEnumeratorEvaluate(t *testing.T) {
	tech := VersionEnumerator{}
	g := &model.Graph{}
	if tech.Evaluate(g) {
		t.Fatalf("expected irrelevance on empty snapshot")
	}
	g = &model.Graph{EvidenceNodes: 1, HypothesisNodes: 1, TechniqueNodes: 1, HasSupportsEdge: true, HasEnablesEdge: true}
	if !tech.Evaluate(g) {
		t.Fatalf("expected relevance for populated snapshot")
	}
}
