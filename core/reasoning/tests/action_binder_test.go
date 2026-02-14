package tests

import (
	"testing"

	"vantage/core/reasoning"
	"vantage/core/state"
)

func TestMatchPatternsRequiresNodesAndEdges(t *testing.T) {
	g := reasoning.NewGraph()
	g.UpsertNode(&reasoning.Node{ID: "ev-1", Type: reasoning.NodeTypeEvidence, Label: "evidence"})
	g.UpsertNode(&reasoning.Node{ID: "hyp-1", Type: reasoning.NodeTypeHypothesis, Label: "hyp"})
	if err := g.AddEdge(&reasoning.Edge{From: "ev-1", To: "hyp-1", Type: reasoning.EdgeTypeSupports}); err != nil {
		t.Fatalf("add edge: %v", err)
	}

	ok := reasoning.MatchPatterns(g, []reasoning.GraphPattern{{
		RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence, reasoning.NodeTypeHypothesis},
		RequiredEdges:     []reasoning.EdgeType{reasoning.EdgeTypeSupports},
	}})
	if !ok {
		t.Fatalf("expected pattern to match")
	}
}

func TestActionClassHypothesisGeneration(t *testing.T) {
	eng := reasoning.NewEngine(nil)
	st, _ := state.New("camp")

	eng.BindActionClasses([]reasoning.ActionClass{{
		ID:            "AC-TEST",
		Name:          "Test",
		Phase:         state.PhaseRecon,
		Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}},
	}})
	if err := eng.IngestEvidence(reasoning.EvidenceEvent{TechniqueID: "T-1", Target: "host", Success: true}); err != nil {
		t.Fatalf("ingest: %v", err)
	}
	eng.ConfigureCycle(reasoning.CycleConfig{Target: "host", AllowedTechniques: []string{"T-1"}, Executor: &executorStub{err: nil}})
	// set engine phase context
	_, _ = eng.RunCycle(st)

	h := eng.GenerateHypotheses()
	found := false
	for _, hyp := range h {
		if hyp.ActionClassID == "AC-TEST" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected action-class derived hypothesis")
	}
}

func TestApplyActionMutatesGraph(t *testing.T) {
	binder := reasoning.NewDefaultActionBinder()
	g := reasoning.NewGraph()
	ac := reasoning.ActionClass{ID: "AC-APPLY", ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}, ProducesEdges: []reasoning.EdgeType{reasoning.EdgeTypeSupports}}
	if err := binder.ApplyAction(g, ac, reasoning.EvidenceEvent{TechniqueID: "T-1", Target: "host", Success: true}); err != nil {
		t.Fatalf("apply action: %v", err)
	}
	if len(g.NodesByType(reasoning.NodeTypeHypothesis)) == 0 {
		t.Fatalf("expected produced hypothesis node")
	}
	if !g.HasEdgeType(reasoning.EdgeTypeSupports) {
		t.Fatalf("expected produced supports edge")
	}
}

func TestPhaseRestrictionEnforced(t *testing.T) {
	binder := reasoning.NewDefaultActionBinder()
	binder.BindActionClasses([]reasoning.ActionClass{{
		ID:            "AC-PHASE",
		Name:          "phase restricted",
		Phase:         state.PhaseObjective,
		Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}},
	}})
	g := reasoning.NewGraph()
	g.UpsertNode(&reasoning.Node{ID: "ev-1", Type: reasoning.NodeTypeEvidence})
	st, _ := state.New("camp")

	h, err := binder.MatchAndGenerate(g, st)
	if err != nil {
		t.Fatalf("match: %v", err)
	}
	if len(h) != 0 {
		t.Fatalf("expected no hypotheses for mismatched phase")
	}
}
