package tests

import (
	"testing"

	"vantage/core/reasoning"
	"vantage/core/state"
)

func TestExpandAttackPathsLinear(t *testing.T) {
	eng := reasoning.NewEngine(nil)
	eng.BindActionClasses([]reasoning.ActionClass{
		{ID: "AC-1", Name: "one", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}, ImpactWeight: 1.0, RiskWeight: 0.2},
		{ID: "AC-2", Name: "two", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeTechnique}, ImpactWeight: 1.1, RiskWeight: 0.2},
	})
	eng.Graph().UpsertNode(&reasoning.Node{ID: "ev-1", Type: reasoning.NodeTypeEvidence, Label: "seed"})

	st, _ := state.New("campaign-linear")
	paths, err := eng.ExpandAttackPaths(st)
	if err != nil {
		t.Fatalf("expand attack paths: %v", err)
	}
	if len(paths) == 0 {
		t.Fatalf("expected at least one path")
	}
	if !paths[0].Valid {
		t.Fatalf("expected top path to be valid")
	}
}

func TestExpandAttackPathsBranching(t *testing.T) {
	eng := reasoning.NewEngine(nil)
	eng.BindActionClasses([]reasoning.ActionClass{
		{ID: "AC-A", Name: "branch-a", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeTechnique}, ImpactWeight: 1.2, RiskWeight: 0.1},
		{ID: "AC-B", Name: "branch-b", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeTechnique}, ImpactWeight: 0.6, RiskWeight: 0.1},
	})
	eng.Graph().UpsertNode(&reasoning.Node{ID: "ev-1", Type: reasoning.NodeTypeEvidence, Label: "seed"})

	st, _ := state.New("campaign-branch")
	paths, err := eng.ExpandAttackPaths(st)
	if err != nil {
		t.Fatalf("expand attack paths: %v", err)
	}
	if len(paths) < 2 {
		t.Fatalf("expected branching paths, got %d", len(paths))
	}
	if paths[0].Score < paths[1].Score {
		t.Fatalf("expected descending score sort")
	}
}

func TestExpandAttackPathsPruningByRiskAndDepth(t *testing.T) {
	eng := reasoning.NewEngine(nil)
	eng.ConfigureAttackPathExpansion(reasoning.AttackPathConfig{
		MaxDepth:           1,
		RiskThreshold:      0.3,
		DepthPenalty:       0.1,
		ConfidenceWeight:   0.2,
		StartNodeTypes:     []reasoning.NodeType{reasoning.NodeTypeEvidence},
		ObjectiveNodeTypes: []reasoning.NodeType{reasoning.NodeTypeTechnique},
	})
	eng.BindActionClasses([]reasoning.ActionClass{
		{ID: "AC-RISKY", Name: "risky", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeTechnique}, ImpactWeight: 2.0, RiskWeight: 0.9},
		{ID: "AC-SAFE", Name: "safe", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeTechnique}, ImpactWeight: 0.7, RiskWeight: 0.1},
	})
	eng.Graph().UpsertNode(&reasoning.Node{ID: "ev-1", Type: reasoning.NodeTypeEvidence, Label: "seed"})
	st, _ := state.New("campaign-prune")

	paths, err := eng.ExpandAttackPaths(st)
	if err != nil {
		t.Fatalf("expand attack paths: %v", err)
	}
	if len(paths) == 0 {
		t.Fatalf("expected at least one surviving path")
	}
	for _, p := range paths {
		if p.Risk > 0.3 {
			t.Fatalf("path with risk %.2f should have been pruned", p.Risk)
		}
		if len(p.Steps) > 1 {
			t.Fatalf("path depth %d exceeded max depth", len(p.Steps))
		}
	}
}

func TestExpandAttackPathsObjectiveReachability(t *testing.T) {
	eng := reasoning.NewEngine(nil)
	eng.BindActionClasses([]reasoning.ActionClass{{
		ID: "AC-OBJ", Name: "objective", Phase: state.PhaseRecon,
		Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}},
		ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeAttackPath},
		ImpactWeight:  1.0,
		RiskWeight:    0.1,
	}})
	eng.ConfigureAttackPathExpansion(reasoning.AttackPathConfig{
		MaxDepth:           2,
		RiskThreshold:      1,
		DepthPenalty:       0.1,
		ConfidenceWeight:   0.2,
		StartNodeTypes:     []reasoning.NodeType{reasoning.NodeTypeEvidence},
		ObjectiveNodeTypes: []reasoning.NodeType{reasoning.NodeTypeAttackPath},
	})
	eng.Graph().UpsertNode(&reasoning.Node{ID: "ev-1", Type: reasoning.NodeTypeEvidence, Label: "seed"})
	st, _ := state.New("campaign-obj")

	paths, err := eng.ExpandAttackPaths(st)
	if err != nil {
		t.Fatalf("expand attack paths: %v", err)
	}
	if len(paths) == 0 {
		t.Fatalf("expected objective path")
	}
	if paths[0].Objective != reasoning.NodeTypeAttackPath {
		t.Fatalf("expected objective type %s, got %s", reasoning.NodeTypeAttackPath, paths[0].Objective)
	}
}
