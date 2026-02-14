package tests

import (
	"testing"

	"vantage/core/reasoning"
	"vantage/core/state"
)

func TestRicherEnvironmentDrivesDeeperAndMorePaths(t *testing.T) {
	minimal := buildSeededEngine(reasoning.SeedScenarioMinimal)
	rich := buildSeededEngine(reasoning.SeedScenarioRich)
	st, _ := state.New("sim")

	minimal.ConfigureAttackPathExpansion(reasoning.AttackPathConfig{MaxDepth: 4, BeamWidth: 25, RiskThreshold: 4, StartNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence, reasoning.NodeTypeHypothesis}, ObjectiveNodeTypes: []reasoning.NodeType{reasoning.NodeTypeDataExposure}, ROEPolicy: func(reasoning.ActionClass, *reasoning.Graph, *state.State) bool { return true }})
	rich.ConfigureAttackPathExpansion(reasoning.AttackPathConfig{MaxDepth: 6, BeamWidth: 25, RiskThreshold: 4, StartNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence, reasoning.NodeTypeHypothesis}, ObjectiveNodeTypes: []reasoning.NodeType{reasoning.NodeTypeDataExposure}, ROEPolicy: func(reasoning.ActionClass, *reasoning.Graph, *state.State) bool { return true }})

	minPaths, _ := minimal.ExpandAttackPaths(st)
	richPaths, _ := rich.ExpandAttackPaths(st)
	if len(richPaths) <= len(minPaths) {
		t.Fatalf("expected richer environment to scale path count min=%d rich=%d", len(minPaths), len(richPaths))
	}
	if maxDepth(richPaths) < maxDepth(minPaths) {
		t.Fatalf("expected richer environment to support deeper paths")
	}
}

func TestObjectiveAppearsOnlyWhenGraphSupportsIt(t *testing.T) {
	eng := buildSeededEngine(reasoning.SeedScenarioMinimal)
	st, _ := state.New("sim2")
	eng.ConfigureAttackPathExpansion(reasoning.AttackPathConfig{MaxDepth: 5, BeamWidth: 20, RiskThreshold: 4, StartNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence, reasoning.NodeTypeHypothesis}, ObjectiveNodeTypes: []reasoning.NodeType{reasoning.NodeTypePrivEsc}, ROEPolicy: func(reasoning.ActionClass, *reasoning.Graph, *state.State) bool { return true }})

	paths, _ := eng.ExpandAttackPaths(st)
	for _, p := range paths {
		if p.Objective == reasoning.NodeTypePrivEsc {
			t.Fatalf("minimal graph should not support priv-esc objective")
		}
	}

	rich := buildSeededEngine(reasoning.SeedScenarioRich)
	rich.ConfigureAttackPathExpansion(reasoning.AttackPathConfig{MaxDepth: 5, BeamWidth: 20, RiskThreshold: 4, StartNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence, reasoning.NodeTypeHypothesis}, ObjectiveNodeTypes: []reasoning.NodeType{reasoning.NodeTypePrivEsc}, ROEPolicy: func(reasoning.ActionClass, *reasoning.Graph, *state.State) bool { return true }})
	richPaths, _ := rich.ExpandAttackPaths(st)
	found := false
	for _, p := range richPaths {
		if p.Objective == reasoning.NodeTypePrivEsc {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected rich graph to expose objective path")
	}
}

func buildSeededEngine(s reasoning.SeedScenario) *reasoning.Engine {
	eng := reasoning.NewEngine(nil)
	reasoning.SeedSyntheticEnvironment(eng.Graph(), s)
	classes := []reasoning.ActionClass{
		{ID: "AC-R1", Name: "web recon", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}, RiskWeight: 0.1, ConfidenceBoost: 0.2},
		{ID: "AC-R2", Name: "cloud enum", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeTechnique}, RiskWeight: 0.1, ConfidenceBoost: 0.2},
		{ID: "AC-O1", Name: "exposure objective", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeTechnique}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeDataExposure}, RiskWeight: 0.2, ConfidenceBoost: 0.25},
	}
	if s == reasoning.SeedScenarioRich {
		classes = append(classes,
			reasoning.ActionClass{ID: "AC-R3", Name: "credential pivot", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence, reasoning.NodeTypeHypothesis}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeTechnique}, ProducesEdges: []reasoning.EdgeType{reasoning.EdgeTypeEnables}, RiskWeight: 0.15, ConfidenceBoost: 0.22},
			reasoning.ActionClass{ID: "AC-O2", Name: "priv esc objective", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeTechnique, reasoning.NodeTypeHypothesis}, RequiredEdges: []reasoning.EdgeType{reasoning.EdgeTypeEnables}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypePrivEsc}, RiskWeight: 0.2, ConfidenceBoost: 0.3},
		)
	}
	eng.BindActionClasses(classes)
	return eng
}

func maxDepth(paths []reasoning.AttackPath) int {
	m := 0
	for _, p := range paths {
		if len(p.Steps) > m {
			m = len(p.Steps)
		}
	}
	return m
}
