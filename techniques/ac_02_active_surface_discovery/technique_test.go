package ac_02_active_surface_discovery

import (
	"testing"

	"vantage/techniques/model"
)

func TestSurfaceProbeActionClassID(t *testing.T) {
	tech := SurfaceProbe{}
	if tech.ActionClassID() != "AC-02" {
		t.Fatalf("unexpected action class: %s", tech.ActionClassID())
	}
}

func TestSurfaceProbeEvaluate(t *testing.T) {
	tech := SurfaceProbe{}
	g := &model.Graph{}
	if tech.Evaluate(g) {
		t.Fatalf("expected irrelevance on empty snapshot")
	}
	g = &model.Graph{EvidenceNodes: 1, HypothesisNodes: 1, TechniqueNodes: 1, HasSupportsEdge: true, HasEnablesEdge: true}
	if !tech.Evaluate(g) {
		t.Fatalf("expected relevance for populated snapshot")
	}
}
