package tests

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"testing"

	"vantage/core/reasoning"
	"vantage/core/state"
)

func TestSimulationScenarios(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		maxDepth   int
		setupGraph func(*testing.T, *reasoning.Engine)
		classes    []reasoning.ActionClass
		assertions func(*testing.T, []reasoning.AttackPath, int)
	}{
		{
			name:     "MinimalExposure",
			maxDepth: 3,
			setupGraph: func(t *testing.T, eng *reasoning.Engine) {
				t.Helper()
				seedGraph(t, eng, []reasoning.Node{
					{ID: "min-ev", Type: reasoning.NodeTypeEvidence, Label: "external footprint"},
					{ID: "min-hyp", Type: reasoning.NodeTypeHypothesis, Label: "surface hypothesis"},
				}, []reasoning.Edge{{From: "min-ev", To: "min-hyp", Type: reasoning.EdgeTypeSupports, Weight: 1}})
			},
			classes: []reasoning.ActionClass{
				{ID: "AC-01", Name: "Recon", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}, RiskWeight: 0.2, ImpactWeight: 0.4, ConfidenceBoost: 0.1},
				{ID: "AC-02", Name: "Validation", Phase: state.PhaseInitialAccess, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeAttackPath}, RiskWeight: 0.3, ImpactWeight: 0.7, ConfidenceBoost: 0.2},
			},
			assertions: func(t *testing.T, paths []reasoning.AttackPath, maxDepth int) {
				t.Helper()
				if len(paths) == 0 {
					t.Fatalf("expected non-empty paths for MinimalExposure")
				}
				if maxDepth > 3 {
					t.Fatalf("expected depth <= 3, got %d", maxDepth)
				}
				if hasActionClass(paths, "AC-13") || hasActionClass(paths, "AC-15") {
					t.Fatalf("expected no AC-13 or AC-15 in MinimalExposure paths")
				}
			},
		},
		{
			name:     "CredentialLeak",
			maxDepth: 5,
			setupGraph: func(t *testing.T, eng *reasoning.Engine) {
				t.Helper()
				seedGraph(t, eng, []reasoning.Node{
					{ID: "cred-ev", Type: reasoning.NodeTypeEvidence, Label: "credential residue"},
					{ID: "cred-tech", Type: reasoning.NodeTypeTechnique, Label: "legacy auth"},
					{ID: "cred-hyp", Type: reasoning.NodeTypeHypothesis, Label: "credential hypothesis"},
				}, []reasoning.Edge{
					{From: "cred-ev", To: "cred-hyp", Type: reasoning.EdgeTypeSupports, Weight: 1},
					{From: "cred-hyp", To: "cred-tech", Type: reasoning.EdgeTypeEnables, Weight: 1},
				})
			},
			classes: []reasoning.ActionClass{
				{ID: "AC-03", Name: "Collect", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeEvidence}, RiskWeight: 0.2, ImpactWeight: 0.2, ConfidenceBoost: 0.05},
				{ID: "AC-04", Name: "Correlate", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}, ProducesEdges: []reasoning.EdgeType{reasoning.EdgeTypeSupports}, RiskWeight: 0.3, ImpactWeight: 0.3, ConfidenceBoost: 0.1},
				{ID: "AC-05", Name: "Stage Access", Phase: state.PhaseInitialAccess, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeTechnique}, ProducesEdges: []reasoning.EdgeType{reasoning.EdgeTypeEnables}, RiskWeight: 0.4, ImpactWeight: 0.4, ConfidenceBoost: 0.1},
				{ID: "AC-08", Name: "Credential Abuse", Phase: state.PhaseInitialAccess, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence, reasoning.NodeTypeHypothesis, reasoning.NodeTypeTechnique}, RequiredEdges: []reasoning.EdgeType{reasoning.EdgeTypeSupports, reasoning.EdgeTypeEnables}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeAttackPath}, RiskWeight: 0.8, ImpactWeight: 0.9, ConfidenceBoost: 0.2},
				{ID: "AC-09", Name: "Token Replay", Phase: state.PhaseInitialAccess, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence, reasoning.NodeTypeHypothesis, reasoning.NodeTypeTechnique}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeAttackPath}, RiskWeight: 0.7, ImpactWeight: 0.85, ConfidenceBoost: 0.2},
			},
			assertions: func(t *testing.T, paths []reasoning.AttackPath, maxDepth int) {
				t.Helper()
				if len(paths) == 0 {
					t.Fatalf("expected non-empty paths for CredentialLeak")
				}
				if maxDepth < 4 {
					t.Fatalf("expected depth >= 4, got %d", maxDepth)
				}
				if !(hasActionClass(paths, "AC-08") || hasActionClass(paths, "AC-09")) {
					t.Fatalf("expected AC-08 or AC-09 in CredentialLeak paths")
				}
			},
		},
		{
			name:     "InternalAccess",
			maxDepth: 4,
			setupGraph: func(t *testing.T, eng *reasoning.Engine) {
				t.Helper()
				seedGraph(t, eng, []reasoning.Node{
					{ID: "int-ev", Type: reasoning.NodeTypeEvidence, Label: "internal network map"},
					{ID: "int-hyp", Type: reasoning.NodeTypeHypothesis, Label: "pivot candidate"},
				}, []reasoning.Edge{{From: "int-ev", To: "int-hyp", Type: reasoning.EdgeTypeSupports, Weight: 1}})
			},
			classes: []reasoning.ActionClass{
				{ID: "AC-10", Name: "Pivot Discovery", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeTechnique}, RiskWeight: 0.3, ImpactWeight: 0.5, ConfidenceBoost: 0.1},
				{ID: "AC-11", Name: "Internal Access", Phase: state.PhaseInitialAccess, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence, reasoning.NodeTypeTechnique}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeAttackPath}, RiskWeight: 0.6, ImpactWeight: 0.8, ConfidenceBoost: 0.15},
			},
			assertions: func(t *testing.T, paths []reasoning.AttackPath, _ int) {
				t.Helper()
				if len(paths) == 0 {
					t.Fatalf("expected non-empty paths for InternalAccess")
				}
				if !hasActionClass(paths, "AC-11") {
					t.Fatalf("expected AC-11 in InternalAccess paths")
				}
			},
		},
		{
			name:     "HighValueAsset",
			maxDepth: 5,
			setupGraph: func(t *testing.T, eng *reasoning.Engine) {
				t.Helper()
				seedGraph(t, eng, []reasoning.Node{
					{ID: "hva-ev", Type: reasoning.NodeTypeEvidence, Label: "critical asset telemetry"},
					{ID: "hva-hyp", Type: reasoning.NodeTypeHypothesis, Label: "high value exposure"},
					{ID: "hva-tech", Type: reasoning.NodeTypeTechnique, Label: "sensitive control plane"},
				}, []reasoning.Edge{{From: "hva-ev", To: "hva-hyp", Type: reasoning.EdgeTypeSupports, Weight: 1}})
			},
			classes: []reasoning.ActionClass{
				{ID: "AC-12", Name: "Asset Validation", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence, reasoning.NodeTypeHypothesis}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeTechnique}, RiskWeight: 0.4, ImpactWeight: 0.7, ConfidenceBoost: 0.1},
				{ID: "AC-13", Name: "Data Impact", Phase: state.PhaseInitialAccess, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence, reasoning.NodeTypeTechnique}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeAttackPath}, RiskWeight: 0.9, ImpactWeight: 1.0, ConfidenceBoost: 0.25},
				{ID: "AC-14", Name: "Control Weakening", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}, RiskWeight: 0.5, ImpactWeight: 0.6, ConfidenceBoost: 0.1},
				{ID: "AC-15", Name: "Operational Disruption", Phase: state.PhaseInitialAccess, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeHypothesis, reasoning.NodeTypeTechnique}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeAttackPath}, RiskWeight: 1.0, ImpactWeight: 1.2, ConfidenceBoost: 0.3},
			},
			assertions: func(t *testing.T, paths []reasoning.AttackPath, _ int) {
				t.Helper()
				if len(paths) == 0 {
					t.Fatalf("expected non-empty paths for HighValueAsset")
				}
				if !(hasActionClass(paths, "AC-13") || hasActionClass(paths, "AC-15")) {
					t.Fatalf("expected AC-13 or AC-15 in HighValueAsset paths")
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			eng := reasoning.NewEngine(nil)
			eng.BindActionClasses(tc.classes)
			tc.setupGraph(t, eng)

			st, err := state.New("simulation-" + strings.ToLower(tc.name))
			if err != nil {
				t.Fatalf("new state: %v", err)
			}

			eng.ConfigureAttackPathExpansion(reasoning.AttackPathConfig{
				MaxDepth:           tc.maxDepth,
				RiskThreshold:      10,
				DepthPenalty:       0.1,
				ConfidenceWeight:   0.25,
				StartNodeTypes:     []reasoning.NodeType{reasoning.NodeTypeEvidence, reasoning.NodeTypeHypothesis, reasoning.NodeTypeTechnique},
				ObjectiveNodeTypes: []reasoning.NodeType{reasoning.NodeTypeAttackPath},
				ROEPolicy:          func(reasoning.ActionClass, *reasoning.Graph, *state.State) bool { return true },
			})

			paths, err := eng.ExpandAttackPaths(st)
			if err != nil {
				t.Fatalf("expand attack paths: %v", err)
			}

			sort.Slice(paths, func(i, j int) bool {
				if paths[i].Score == paths[j].Score {
					return len(paths[i].Steps) < len(paths[j].Steps)
				}
				return paths[i].Score > paths[j].Score
			})

			maxDepth := 0
			scoreMin := math.MaxFloat64
			scoreMax := -math.MaxFloat64
			for _, p := range paths {
				if len(p.Steps) > maxDepth {
					maxDepth = len(p.Steps)
				}
				if p.Score < scoreMin {
					scoreMin = p.Score
				}
				if p.Score > scoreMax {
					scoreMax = p.Score
				}
			}
			if len(paths) == 0 {
				scoreMin = 0
				scoreMax = 0
			}

			t.Logf("Total path count: %d", len(paths))
			t.Logf("Max depth: %d", maxDepth)
			t.Logf("Score range: %.3f - %.3f", scoreMin, scoreMax)
			for i, seq := range topActionSequences(paths, 3) {
				t.Logf("Top %d path actions: %s", i+1, seq)
			}

			tc.assertions(t, paths, maxDepth)
		})
	}
}

func seedGraph(t *testing.T, eng *reasoning.Engine, nodes []reasoning.Node, edges []reasoning.Edge) {
	t.Helper()
	for i := range nodes {
		n := nodes[i]
		eng.Graph().UpsertNode(&n)
	}
	for i := range edges {
		e := edges[i]
		if err := eng.Graph().AddEdge(&e); err != nil {
			t.Fatalf("add edge %s->%s: %v", e.From, e.To, err)
		}
	}
}

func topActionSequences(paths []reasoning.AttackPath, limit int) []string {
	if limit > len(paths) {
		limit = len(paths)
	}
	out := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		parts := make([]string, 0, len(paths[i].Steps))
		for _, step := range paths[i].Steps {
			parts = append(parts, step.ActionClassID)
		}
		out = append(out, fmt.Sprintf("[%s]", strings.Join(parts, " -> ")))
	}
	return out
}

func hasActionClass(paths []reasoning.AttackPath, id string) bool {
	for _, p := range paths {
		for _, step := range p.Steps {
			if step.ActionClassID == id {
				return true
			}
		}
	}
	return false
}
