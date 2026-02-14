package tests

import (
	"context"
	"errors"
	"strings"
	"testing"

	"vantage/core/evidence"
	"vantage/core/reasoning"
	"vantage/core/state"
)

type fixedExpander struct {
	hypotheses []reasoning.Hypothesis
}

func (f fixedExpander) Expand(_ *reasoning.Graph, _ *state.State) ([]reasoning.Hypothesis, error) {
	return f.hypotheses, nil
}

type failingExpander struct{}

func (failingExpander) Expand(_ *reasoning.Graph, _ *state.State) ([]reasoning.Hypothesis, error) {
	return nil, errors.New("ai unavailable")
}

type executorStub struct {
	artifact *evidence.Artifact
	err      error
	calls    int
	lastTech string
	lastTgt  string
}

func (s *executorStub) Run(_ context.Context, techniqueID string, target string) (*evidence.Artifact, error) {
	s.calls++
	s.lastTech = techniqueID
	s.lastTgt = target
	return s.artifact, s.err
}

func TestReasoningCyclePlansHighestScore(t *testing.T) {
	re := reasoning.NewEngine(nil)
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
	re := reasoning.NewEngine(nil)
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

func TestReasoningEngineMergesDeterministicAndAIHypotheses(t *testing.T) {
	re := reasoning.NewEngine(fixedExpander{hypotheses: []reasoning.Hypothesis{{
		ID:         "hyp-ai-1",
		Statement:  "ai generated",
		Confidence: 0.61,
	}}})
	re.RegisterTechniqueEffect(reasoning.TechniqueEffect{TechniqueID: "T-1", Impact: 0.8, Risk: 0.2, Stealth: 0.7})
	_ = re.IngestEvidence(reasoning.EvidenceEvent{TechniqueID: "T-1", Target: "host-1", Success: true})

	_, err := re.PlanNextAction(reasoning.PlannerQuery{Target: "host-1", AllowedTechniques: []string{"T-1"}})
	if err != nil {
		t.Fatalf("plan next action: %v", err)
	}

	hypNodes := re.Graph().NodesByType(reasoning.NodeTypeHypothesis)
	if len(hypNodes) != 2 {
		t.Fatalf("expected deterministic + ai hypotheses, got %d", len(hypNodes))
	}

	foundAI := false
	for _, n := range hypNodes {
		if n.ID == "hyp-ai-1" {
			foundAI = true
			break
		}
	}
	if !foundAI {
		t.Fatalf("expected merged ai hypothesis node")
	}
}

func TestReasoningEngineAIFailureDoesNotBreakCycle(t *testing.T) {
	re := reasoning.NewEngine(failingExpander{})
	re.RegisterTechniqueEffect(reasoning.TechniqueEffect{TechniqueID: "T-1", Impact: 0.8, Risk: 0.2, Stealth: 0.7})
	_ = re.IngestEvidence(reasoning.EvidenceEvent{TechniqueID: "T-1", Target: "host-1", Success: true})

	decision, err := re.PlanNextAction(reasoning.PlannerQuery{Target: "host-1", AllowedTechniques: []string{"T-1"}})
	if err != nil {
		t.Fatalf("plan next action should succeed when ai expander fails: %v", err)
	}
	if decision.Selected.TechniqueID != "T-1" {
		t.Fatalf("expected deterministic selection, got %s", decision.Selected.TechniqueID)
	}
}

func TestRunCycleInvokesExecutorAndIngestsEvidence(t *testing.T) {
	re := reasoning.NewEngine(nil)
	re.RegisterTechniqueEffect(reasoning.TechniqueEffect{TechniqueID: "T-1", Impact: 0.9, Risk: 0.1, Stealth: 0.8})

	s := &executorStub{artifact: &evidence.Artifact{TechniqueID: "T-1", Target: "host-1", Success: true}}
	re.ConfigureCycle(reasoning.CycleConfig{Target: "host-1", AllowedTechniques: []string{"T-1"}, Executor: s})

	st, err := state.New("campaign-1")
	if err != nil {
		t.Fatalf("state new: %v", err)
	}
	decision, err := re.RunCycle(st)
	if err != nil {
		t.Fatalf("run cycle: %v", err)
	}
	if decision.Selected.TechniqueID != "T-1" {
		t.Fatalf("unexpected selected technique: %s", decision.Selected.TechniqueID)
	}
	if s.calls != 1 {
		t.Fatalf("expected executor call count 1, got %d", s.calls)
	}
	if len(re.Graph().NodesByType(reasoning.NodeTypeEvidence)) == 0 {
		t.Fatalf("expected evidence node to be ingested")
	}
}

func TestRunCycleReturnsDecisionAndExecutorError(t *testing.T) {
	re := reasoning.NewEngine(nil)
	re.RegisterTechniqueEffect(reasoning.TechniqueEffect{TechniqueID: "T-1", Impact: 0.9, Risk: 0.1, Stealth: 0.8})

	s := &executorStub{err: errors.New("execution failed")}
	re.ConfigureCycle(reasoning.CycleConfig{Target: "host-1", AllowedTechniques: []string{"T-1"}, Executor: s})

	st, err := state.New("campaign-2")
	if err != nil {
		t.Fatalf("state new: %v", err)
	}
	decision, runErr := re.RunCycle(st)
	if runErr == nil {
		t.Fatalf("expected run cycle error")
	}
	if decision == nil {
		t.Fatalf("expected decision even on executor error")
	}
	if s.calls != 1 {
		t.Fatalf("expected executor call count 1, got %d", s.calls)
	}
}
