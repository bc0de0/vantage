package ac_15_external_execution

import (
	"testing"

	"vantage/techniques/model"
)

func TestExternalExecutionCoordinatorActionClassID(t *testing.T) {
	tech := ExternalExecutionCoordinator{}
	if tech.ActionClassID() != "AC-15" {
		t.Fatalf("unexpected action class: %s", tech.ActionClassID())
	}
}

func TestExternalExecutionCoordinatorEvaluate(t *testing.T) {
	tech := ExternalExecutionCoordinator{}
	g := &model.Graph{}
	if tech.Evaluate(g) {
		t.Fatalf("expected irrelevance on empty snapshot")
	}
	g = &model.Graph{EvidenceNodes: 1, HypothesisNodes: 1, TechniqueNodes: 1, HasSupportsEdge: true, HasEnablesEdge: true}
	if !tech.Evaluate(g) {
		t.Fatalf("expected relevance for populated snapshot")
	}
}
