package tests

import (
	"testing"

	"vantage/core/reasoning"
	"vantage/core/state"
)

func TestPlanCampaignGoalReachabilityAndOrdering(t *testing.T) {
	eng := reasoning.NewEngine(nil)
	eng.BindActionClasses([]reasoning.ActionClass{
		{ID: "AC-R", Name: "recon", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}, RiskWeight: 0.1, ConfidenceBoost: 0.2},
		{ID: "AC-D", Name: "data", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeDataExposure}, RiskWeight: 0.2, ConfidenceBoost: 0.3},
		{ID: "AC-P", Name: "priv", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypePrivEsc}, RiskWeight: 0.25, ConfidenceBoost: 0.25},
	})
	eng.Graph().UpsertNode(&reasoning.Node{ID: "seed", Type: reasoning.NodeTypeEvidence, Label: "seed"})

	campaigns, err := eng.PlanCampaign(reasoning.NodeTypeDataExposure, reasoning.CampaignOptions{MaxDepth: 4, RiskTolerance: 1.0, ConfidenceThreshold: 0.4, BeamWidth: 6, TopN: 5})
	if err != nil {
		t.Fatalf("plan campaign: %v", err)
	}
	if len(campaigns) == 0 {
		t.Fatalf("expected campaigns")
	}
	for i := range campaigns {
		if campaigns[i].Objective != reasoning.NodeTypeDataExposure {
			t.Fatalf("objective mismatch")
		}
		last := campaigns[i].Steps[len(campaigns[i].Steps)-1]
		if last.ActionClassID != "AC-D" {
			t.Fatalf("campaign does not reach objective action")
		}
		if i > 0 && campaigns[i-1].Score < campaigns[i].Score {
			t.Fatalf("scores are not descending")
		}
	}
}

func TestPlanCampaignBeamPruningAndNoExplosion(t *testing.T) {
	eng := reasoning.NewEngine(nil)
	classes := make([]reasoning.ActionClass, 0, 13)
	for i := 0; i < 12; i++ {
		classes = append(classes, reasoning.ActionClass{ID: "AC-PRE-" + string(rune('A'+i)), Name: "prep", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}, RiskWeight: 0.05, ConfidenceBoost: 0.15})
	}
	classes = append(classes, reasoning.ActionClass{ID: "AC-OBJ", Name: "objective", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeDataExposure}, RiskWeight: 0.05, ConfidenceBoost: 0.3})
	eng.BindActionClasses(classes)
	eng.Graph().UpsertNode(&reasoning.Node{ID: "seed", Type: reasoning.NodeTypeEvidence, Label: "seed"})

	narrow, err := eng.PlanCampaign(reasoning.NodeTypeDataExposure, reasoning.CampaignOptions{MaxDepth: 5, BeamWidth: 2, RiskTolerance: 3, ConfidenceThreshold: 0.2, TopN: 100})
	if err != nil {
		t.Fatalf("narrow plan err: %v", err)
	}
	wide, err := eng.PlanCampaign(reasoning.NodeTypeDataExposure, reasoning.CampaignOptions{MaxDepth: 5, BeamWidth: 8, RiskTolerance: 3, ConfidenceThreshold: 0.2, TopN: 100})
	if err != nil {
		t.Fatalf("wide plan err: %v", err)
	}
	if len(narrow) > len(wide) {
		t.Fatalf("beam pruning broken narrow=%d wide=%d", len(narrow), len(wide))
	}
	if len(wide) > 100 {
		t.Fatalf("unexpected campaign explosion: %d", len(wide))
	}
}
