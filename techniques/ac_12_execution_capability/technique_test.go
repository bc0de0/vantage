package ac_12_execution_capability

import (
	"testing"

	"vantage/techniques/model"
)

func TestExecutionCapabilityValidatorActionClassID(t *testing.T) {
	tech := ExecutionCapabilityValidator{}
	if tech.ActionClassID() != "AC-12" {
		t.Fatalf("unexpected action class: %s", tech.ActionClassID())
	}
}

func TestExecutionCapabilityValidatorEvaluate(t *testing.T) {
	tech := ExecutionCapabilityValidator{}
	g := &model.Graph{}
	if tech.Evaluate(g) {
		t.Fatalf("expected irrelevance on empty snapshot")
	}
	g = &model.Graph{EvidenceNodes: 1, HypothesisNodes: 1, TechniqueNodes: 1, HasSupportsEdge: true, HasEnablesEdge: true}
	if !tech.Evaluate(g) {
		t.Fatalf("expected relevance for populated snapshot")
	}
}
