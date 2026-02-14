package tests

import (
	"strings"
	"testing"

	"vantage/core/reasoning"
)

func TestReasoningCyclePlansHighestScore(t *testing.T) {
	re := reasoning.NewEngine()
	re.RegisterTechniqueEffect(reasoning.TechniqueEffect{TechniqueID: "T-A", Impact: 0.9, Risk: 0.2, Stealth: 0.7})
	re.RegisterTechniqueEffect(reasoning.TechniqueEffect{TechniqueID: "T-B", Impact: 0.5, Risk: 0.1, Stealth: 0.8})

	if err := re.IngestEvidence(reasoning.EvidenceEvent{TechniqueID: "T-A", Target: "host-1", Success: true}); err != nil {
		t.Fatalf("ingest evidence: %v", err)
	}

	decision, err := re.PlanNextAction(reasoning.PlannerQuery{Target: "host-1", AllowedTechniques: []string{"T-A", "T-B"}})
	if err != nil {
		t.Fatalf("plan next action: %v", err)
	}
	if decision.Selected.TechniqueID != "T-A" {
		t.Fatalf("expected T-A to win scoring, got %s", decision.Selected.TechniqueID)
	}
	if len(decision.Ranked) != 2 {
		t.Fatalf("expected 2 ranked actions, got %d", len(decision.Ranked))
	}
}

func TestReasoningGraphDOTIncludesEvidenceAndHypothesis(t *testing.T) {
	re := reasoning.NewEngine()
	re.RegisterTechniqueEffect(reasoning.TechniqueEffect{TechniqueID: "T-X", Impact: 0.8, Risk: 0.3, Stealth: 0.6})
	_ = re.IngestEvidence(reasoning.EvidenceEvent{TechniqueID: "T-X", Target: "target-1", Success: true})
	_, err := re.PlanNextAction(reasoning.PlannerQuery{Target: "target-1", AllowedTechniques: []string{"T-X"}})
	if err != nil {
		t.Fatalf("plan next action: %v", err)
	}
	dot := re.DOT()
	if !strings.Contains(dot, "digraph reasoning") {
		t.Fatalf("expected DOT header")
	}
	if !strings.Contains(dot, "supports") {
		t.Fatalf("expected supports edge in DOT output")
	}
}
