package tests

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"vantage/core/reasoning"
	"vantage/core/state"
)

func TestFullPathExpansionWithNormalizedActionClasses(t *testing.T) {
	// This integration test validates end-to-end deterministic path expansion using
	// the normalized action-class corpus with the default binder/planner/scoring stack
	// and no AI hypothesis expander (NewEngine(nil)).
	eng := reasoning.NewEngine(nil)

	classesDir := resolveActionClassesDir(t)
	classes, err := reasoning.LoadActionClassesFromDir(classesDir)
	if err != nil {
		t.Fatalf("load action classes: %v", err)
	}
	if len(classes) == 0 {
		t.Fatalf("expected normalized action classes to load")
	}

	// Use the full normalized set but tune a small subset so objective reachability,
	// depth expansion, pruning, and scoring are all observable in a single test.
	edited := make([]reasoning.ActionClass, len(classes))
	copy(edited, classes)
	for i := range edited {
		switch edited[i].ID {
		case "AC-01":
			edited[i].ImpactWeight = 1.2
			edited[i].RiskWeight = 0.1
		case "AC-02":
			edited[i].ImpactWeight = 0.7
			edited[i].RiskWeight = 0.6
		case "AC-07":
			edited[i].Preconditions = []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence, reasoning.NodeTypeHypothesis}}}
			edited[i].ProducesNodes = append(edited[i].ProducesNodes, reasoning.NodeTypeTechnique)
			edited[i].ImpactWeight = 0.9
			edited[i].RiskWeight = 0.2
		case "AC-08":
			edited[i].Preconditions = []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence, reasoning.NodeTypeHypothesis}}}
			edited[i].ProducesNodes = append(edited[i].ProducesNodes, reasoning.NodeTypeTechnique)
			edited[i].ImpactWeight = 0.5
			edited[i].RiskWeight = 0.5
		case "AC-09":
			edited[i].ProducesNodes = []reasoning.NodeType{reasoning.NodeTypeAttackPath}
			edited[i].ImpactWeight = 1.1
			edited[i].RiskWeight = 0.2
		}
	}
	eng.BindActionClasses(edited)

	// Minimal graph bootstrap: only the starting target seed node.
	eng.Graph().UpsertNode(&reasoning.Node{ID: "target-seed", Type: reasoning.NodeTypeEvidence, Label: "target"})

	st, err := state.New("campaign-full-path")
	if err != nil {
		t.Fatalf("new state: %v", err)
	}

	// Baseline expansion checks:
	// - returns feasible paths
	// - explores to depth >= 3
	// - keeps phase progression constraints
	// - avoids action-class repetition within a path (via ROE guard)
	eng.ConfigureAttackPathExpansion(reasoning.AttackPathConfig{
		MaxDepth:           4,
		RiskThreshold:      3.0,
		DepthPenalty:       0.1,
		ConfidenceWeight:   0.25,
		StartNodeTypes:     []reasoning.NodeType{reasoning.NodeTypeEvidence},
		ObjectiveNodeTypes: []reasoning.NodeType{reasoning.NodeTypeAttackPath},
		ROEPolicy: func(ac reasoning.ActionClass, g *reasoning.Graph, _ *state.State) bool {
			// Prevent repeating the same action class in a single explored path.
			for _, n := range g.NodesByType(reasoning.NodeTypeEvidence) {
				if n.Label == "simulated "+ac.ID {
					return false
				}
			}
			for _, n := range g.NodesByType(reasoning.NodeTypeHypothesis) {
				if n.Label == "simulated "+ac.ID {
					return false
				}
			}
			for _, n := range g.NodesByType(reasoning.NodeTypeTechnique) {
				if n.Label == "simulated "+ac.ID {
					return false
				}
			}
			return true
		},
	})

	paths, err := eng.ExpandAttackPaths(st)
	if err != nil {
		t.Fatalf("expand attack paths: %v", err)
	}
	// Breadth matters for cognitive validation: strong reasoning should surface
	// multiple distinct options, not a single brittle chain.
	t.Logf("Total paths found: %d", len(paths))
	if len(paths) < 2 {
		t.Fatalf("expected at least two expanded paths for breadth validation, got %d", len(paths))
	}

	maxDepth := 0
	scoreMin := math.MaxFloat64
	scoreMax := -math.MaxFloat64
	scoreSet := map[float64]struct{}{}
	for _, p := range paths {
		if p.Score < scoreMin {
			scoreMin = p.Score
		}
		if p.Score > scoreMax {
			scoreMax = p.Score
		}
		scoreSet[p.Score] = struct{}{}
		if len(p.Steps) > maxDepth {
			maxDepth = len(p.Steps)
		}

		seen := map[string]struct{}{}
		for _, step := range p.Steps {
			if _, exists := seen[step.ActionClassID]; exists {
				t.Fatalf("path contains duplicate action class %q", step.ActionClassID)
			}
			seen[step.ActionClassID] = struct{}{}

			ac, ok := findActionClassByID(edited, step.ActionClassID)
			if !ok {
				t.Fatalf("missing action class %q in bound set", step.ActionClassID)
			}
			if ac.Phase != state.PhaseRecon {
				next, nextOK := state.PhaseRecon.Next()
				if !nextOK || ac.Phase != next {
					t.Fatalf("action class %q phase %s violates ordering from %s", ac.ID, ac.Phase, state.PhaseRecon)
				}
			}
		}
	}
	// Depth highlights whether reasoning is shallow (single-hop) or whether it can
	// build multi-step causal chains toward the objective.
	t.Logf("Max depth: %d", maxDepth)
	if maxDepth < 2 {
		t.Fatalf("expected max path depth >= 2, got %d", maxDepth)
	}
	if maxDepth < 3 {
		t.Fatalf("expected at least one path with length >= 3, got %d", maxDepth)
	}
	// Score richness indicates discriminative reasoning quality: useful engines
	// should rank paths differently instead of collapsing to uniform scores.
	t.Logf("Non-linear score range: %.3f - %.3f (spread %.3f)", scoreMin, scoreMax, scoreMax-scoreMin)
	if scoreMax <= scoreMin {
		t.Fatalf("expected non-zero score variance, got min %.2f max %.2f", scoreMin, scoreMax)
	}
	if len(scoreSet) < 2 {
		t.Fatalf("expected non-uniform path scores, got %d unique score(s)", len(scoreSet))
	}

	// Lowering risk threshold should prune higher-risk branches.
	eng.ConfigureAttackPathExpansion(reasoning.AttackPathConfig{
		MaxDepth:           4,
		RiskThreshold:      1.2,
		DepthPenalty:       0.1,
		ConfidenceWeight:   0.25,
		StartNodeTypes:     []reasoning.NodeType{reasoning.NodeTypeEvidence},
		ObjectiveNodeTypes: []reasoning.NodeType{reasoning.NodeTypeAttackPath},
		ROEPolicy:          func(reasoning.ActionClass, *reasoning.Graph, *state.State) bool { return true },
	})
	prunedByRisk, err := eng.ExpandAttackPaths(st)
	if err != nil {
		t.Fatalf("expand with lower risk threshold: %v", err)
	}
	if len(prunedByRisk) >= len(paths) {
		t.Fatalf("expected fewer paths after lowering risk threshold, got %d (baseline %d)", len(prunedByRisk), len(paths))
	}

	// Lowering max depth should cut off the 3-step objective chain.
	eng.ConfigureAttackPathExpansion(reasoning.AttackPathConfig{
		MaxDepth:           2,
		RiskThreshold:      3.0,
		DepthPenalty:       0.1,
		ConfidenceWeight:   0.25,
		StartNodeTypes:     []reasoning.NodeType{reasoning.NodeTypeEvidence},
		ObjectiveNodeTypes: []reasoning.NodeType{reasoning.NodeTypeAttackPath},
		ROEPolicy:          func(reasoning.ActionClass, *reasoning.Graph, *state.State) bool { return true },
	})
	prunedByDepth, err := eng.ExpandAttackPaths(st)
	if err != nil {
		t.Fatalf("expand with lower max depth: %v", err)
	}
	if len(prunedByDepth) != 0 {
		t.Fatalf("expected no objective paths with max depth=2, got %d", len(prunedByDepth))
	}
}

func findActionClassByID(classes []reasoning.ActionClass, id string) (reasoning.ActionClass, bool) {
	for _, ac := range classes {
		if ac.ID == id {
			return ac, true
		}
	}
	return reasoning.ActionClass{}, false
}

func resolveActionClassesDir(t *testing.T) string {
	t.Helper()
	candidates := []string{"action-classes-normalized", "../../../action-classes-normalized"}
	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			abs, err := filepath.Abs(candidate)
			if err != nil {
				t.Fatalf("abs path for %s: %v", candidate, err)
			}
			return abs
		}
	}
	t.Fatalf("could not locate action-classes-normalized directory")
	return ""
}
