package tests

import (
	"testing"

	"vantage/core/reasoning"
	"vantage/core/state"
	"vantage/techniqueset"
)

func TestTechniqueDiversityAndPathScoreVariance(t *testing.T) {
	techniques := techniqueset.ForActionClass("AC-09")
	if len(techniques) < 5 {
		t.Fatalf("expected at least 5 AC-09 techniques, got %d", len(techniques))
	}

	riskSeen := map[float64]struct{}{}
	impactSeen := map[float64]struct{}{}
	minRisk, maxRisk := 2.0, -1.0
	minImpact, maxImpact := 2.0, -1.0
	for _, tech := range techniques {
		risk := tech.RiskModifier()
		impact := tech.ImpactModifier()
		riskSeen[risk] = struct{}{}
		impactSeen[impact] = struct{}{}
		if risk < minRisk {
			minRisk = risk
		}
		if risk > maxRisk {
			maxRisk = risk
		}
		if impact < minImpact {
			minImpact = impact
		}
		if impact > maxImpact {
			maxImpact = impact
		}
	}
	if len(riskSeen) < 4 || maxRisk-minRisk < 0.5 {
		t.Fatalf("expected meaningful risk distribution, unique=%d range=%.2f", len(riskSeen), maxRisk-minRisk)
	}
	if len(impactSeen) < 4 || maxImpact-minImpact < 0.5 {
		t.Fatalf("expected meaningful impact distribution, unique=%d range=%.2f", len(impactSeen), maxImpact-minImpact)
	}

	baselineVariance := expandAndVariance(t, []reasoning.ActionClass{
		{ID: "AC-01", Name: "A", Phase: state.PhaseRecon, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeAttackPath}, RiskWeight: 0.5, ConfidenceBoost: 0.1},
		{ID: "AC-02", Name: "B", Phase: state.PhaseRecon, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeAttackPath}, RiskWeight: 0.5, ConfidenceBoost: 0.1},
	})
	diverseVariance := expandAndVariance(t, []reasoning.ActionClass{
		{ID: "AC-01", Name: "A", Phase: state.PhaseRecon, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeAttackPath}, RiskWeight: 0.15, ConfidenceBoost: 0.35},
		{ID: "AC-02", Name: "B", Phase: state.PhaseRecon, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeAttackPath}, RiskWeight: 1.3, ConfidenceBoost: 0.0},
	})
	if diverseVariance <= baselineVariance {
		t.Fatalf("expected diversified path setup to increase score variance: baseline=%.6f diverse=%.6f", baselineVariance, diverseVariance)
	}
}

func expandAndVariance(t *testing.T, classes []reasoning.ActionClass) float64 {
	t.Helper()
	eng := reasoning.NewEngine(nil)
	eng.BindActionClasses(classes)
	eng.Graph().UpsertNode(&reasoning.Node{ID: "seed", Type: reasoning.NodeTypeEvidence, Label: "seed"})

	st, err := state.New("variance-check")
	if err != nil {
		t.Fatalf("new state: %v", err)
	}

	eng.ConfigureAttackPathExpansion(reasoning.AttackPathConfig{
		MaxDepth:           1,
		RiskThreshold:      2.0,
		DepthPenalty:       0.1,
		ConfidenceWeight:   0.25,
		StartNodeTypes:     []reasoning.NodeType{reasoning.NodeTypeEvidence},
		ObjectiveNodeTypes: []reasoning.NodeType{reasoning.NodeTypeAttackPath},
		ROEPolicy:          func(reasoning.ActionClass, *reasoning.Graph, *state.State) bool { return true },
	})

	paths, err := eng.ExpandAttackPaths(st)
	if err != nil {
		t.Fatalf("expand attack paths: %v", err)
	}
	if len(paths) < 2 {
		t.Fatalf("expected at least two paths, got %d", len(paths))
	}

	mean := 0.0
	for _, p := range paths {
		mean += p.Score
	}
	mean /= float64(len(paths))
	variance := 0.0
	for _, p := range paths {
		d := p.Score - mean
		variance += d * d
	}
	return variance / float64(len(paths))
}
