package tests

import (
	"testing"

	"vantage/core/reasoning"
)

func TestPlannerRanksOnlyMatchingActionClassTechniques(t *testing.T) {
	g := reasoning.NewGraph()
	g.UpsertNode(&reasoning.Node{ID: "ev-1", Type: reasoning.NodeTypeEvidence, Label: "evidence"})

	p := reasoning.NewPlanner(nil, reasoning.DefaultTechniqueScoreWeights())
	lookup := func(id string) (reasoning.ActionClass, bool) {
		if id == "AC-02" {
			return reasoning.ActionClass{ID: "AC-02", ImpactWeight: 0.6, RiskWeight: 0.4}, true
		}
		return reasoning.ActionClass{}, false
	}

	ranked := p.RankedActionsForHypotheses(g, "target", []reasoning.Hypothesis{{ID: "h1", ActionClassID: "AC-02"}, {ID: "h2", ActionClassID: "AC-99"}}, lookup, 0)
	if len(ranked) == 0 {
		t.Fatalf("expected ranked techniques for AC-02")
	}
	for _, r := range ranked {
		if r.ActionClassID != "AC-02" {
			t.Fatalf("unexpected action class in ranked result: %s", r.ActionClassID)
		}
	}
}
