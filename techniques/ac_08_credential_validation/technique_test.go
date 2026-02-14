package ac_08_credential_validation

import (
	"testing"

	"vantage/techniques/model"
)

func TestCredentialValidatorActionClassID(t *testing.T) {
	tech := CredentialValidator{}
	if tech.ActionClassID() != "AC-08" {
		t.Fatalf("unexpected action class: %s", tech.ActionClassID())
	}
}

func TestCredentialValidatorEvaluate(t *testing.T) {
	tech := CredentialValidator{}
	g := &model.Graph{}
	if tech.Evaluate(g) {
		t.Fatalf("expected irrelevance on empty snapshot")
	}
	g = &model.Graph{EvidenceNodes: 1, HypothesisNodes: 1, TechniqueNodes: 1, HasSupportsEdge: true, HasEnablesEdge: true}
	if !tech.Evaluate(g) {
		t.Fatalf("expected relevance for populated snapshot")
	}
}
