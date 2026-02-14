package tests

import (
	"testing"

	"vantage/core/reasoning"
	"vantage/core/state"
)

func TestPlanCampaignMultipleObjectives(t *testing.T) {
	eng := reasoning.NewEngine(nil)
	eng.BindActionClasses([]reasoning.ActionClass{
		{ID: "AC-R", Name: "recon", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}, RiskWeight: 0.1, ConfidenceBoost: 0.2},
		{ID: "AC-D", Name: "data", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeDataExposure}, RiskWeight: 0.2, ConfidenceBoost: 0.3},
		{ID: "AC-P", Name: "priv", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypePrivEsc}, RiskWeight: 0.25, ConfidenceBoost: 0.25},
	})
	eng.Graph().UpsertNode(&reasoning.Node{ID: "seed", Type: reasoning.NodeTypeEvidence, Label: "seed"})

	for _, objective := range []reasoning.NodeType{reasoning.NodeTypeDataExposure, reasoning.NodeTypePrivEsc} {
		campaigns, err := eng.PlanCampaign(eng.Graph(), objective, reasoning.CampaignOptions{MaxDepth: 4, RiskTolerance: 1.0, ConfidenceThreshold: 0.4, BeamWidth: 6})
		if err != nil {
			t.Fatalf("plan campaign for %s: %v", objective, err)
		}
		if len(campaigns) == 0 {
			t.Fatalf("expected campaigns for objective %s", objective)
		}
		for _, campaign := range campaigns {
			if campaign.Objective != objective {
				t.Fatalf("campaign objective mismatch got=%s want=%s", campaign.Objective, objective)
			}
			last := campaign.Steps[len(campaign.Steps)-1]
			if objective == reasoning.NodeTypeDataExposure && last.ActionClassID != "AC-D" {
				t.Fatalf("expected campaign to terminate in AC-D, got %s", last.ActionClassID)
			}
			if objective == reasoning.NodeTypePrivEsc && last.ActionClassID != "AC-P" {
				t.Fatalf("expected campaign to terminate in AC-P, got %s", last.ActionClassID)
			}
		}
	}
}

func TestPlanCampaignThresholdsAndOrdering(t *testing.T) {
	eng := reasoning.NewEngine(nil)
	eng.BindActionClasses([]reasoning.ActionClass{
		{ID: "AC-A", Name: "safe prep", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}, RiskWeight: 0.1, ConfidenceBoost: 0.25},
		{ID: "AC-B", Name: "risky prep", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}, RiskWeight: 0.8, ConfidenceBoost: 0.05},
		{ID: "AC-C", Name: "objective", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeLateralReachability}, RiskWeight: 0.15, ConfidenceBoost: 0.35},
	})
	eng.Graph().UpsertNode(&reasoning.Node{ID: "seed", Type: reasoning.NodeTypeEvidence, Label: "seed"})

	campaigns, err := eng.PlanCampaign(eng.Graph(), reasoning.NodeTypeLateralReachability, reasoning.CampaignOptions{MaxDepth: 4, RiskTolerance: 0.5, ConfidenceThreshold: 0.6, BeamWidth: 5})
	if err != nil {
		t.Fatalf("plan campaign: %v", err)
	}
	if len(campaigns) == 0 {
		t.Fatalf("expected at least one campaign")
	}
	for i := range campaigns {
		if campaigns[i].Risk > 0.5 {
			t.Fatalf("risk threshold violated: %.3f", campaigns[i].Risk)
		}
		if campaigns[i].Confidence < 0.6 {
			t.Fatalf("confidence threshold violated: %.3f", campaigns[i].Confidence)
		}
		if i > 0 && campaigns[i-1].Score < campaigns[i].Score {
			t.Fatalf("campaign scores not sorted descending")
		}
	}
}

func TestPlanCampaignBeamWidthLimitsExpansion(t *testing.T) {
	eng := reasoning.NewEngine(nil)
	eng.BindActionClasses([]reasoning.ActionClass{
		{ID: "AC-1", Name: "prep1", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}, RiskWeight: 0.1, ConfidenceBoost: 0.2},
		{ID: "AC-2", Name: "prep2", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}, RiskWeight: 0.1, ConfidenceBoost: 0.18},
		{ID: "AC-3", Name: "prep3", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}, RiskWeight: 0.1, ConfidenceBoost: 0.15},
		{ID: "AC-4", Name: "objective", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeDataExposure}, RiskWeight: 0.1, ConfidenceBoost: 0.3},
	})
	eng.Graph().UpsertNode(&reasoning.Node{ID: "seed", Type: reasoning.NodeTypeEvidence, Label: "seed"})

	narrow, err := eng.PlanCampaign(eng.Graph(), reasoning.NodeTypeDataExposure, reasoning.CampaignOptions{MaxDepth: 4, RiskTolerance: 2.0, ConfidenceThreshold: 0.4, BeamWidth: 1})
	if err != nil {
		t.Fatalf("plan campaign narrow beam: %v", err)
	}
	wide, err := eng.PlanCampaign(eng.Graph(), reasoning.NodeTypeDataExposure, reasoning.CampaignOptions{MaxDepth: 4, RiskTolerance: 2.0, ConfidenceThreshold: 0.4, BeamWidth: 4})
	if err != nil {
		t.Fatalf("plan campaign wide beam: %v", err)
	}
	if len(narrow) > len(wide) {
		t.Fatalf("beam width 1 should not produce more campaigns than wider beam: narrow=%d wide=%d", len(narrow), len(wide))
	}
}
