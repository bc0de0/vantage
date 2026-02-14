package tests

import (
	"fmt"
	"testing"
	"time"

	"vantage/core/reasoning"
	"vantage/core/state"
)

func TestPlanCampaignStress1000Techniques(t *testing.T) {
	eng := reasoning.NewEngine(nil)
	reasoning.SeedSyntheticEnvironment(eng.Graph(), reasoning.SeedScenarioRich)
	classes := make([]reasoning.ActionClass, 0, 1000)
	for i := 0; i < 1000; i++ {
		id := fmt.Sprintf("AC-X-%04d", i)
		produces := []reasoning.NodeType{reasoning.NodeTypeHypothesis}
		if i%50 == 0 {
			produces = []reasoning.NodeType{reasoning.NodeTypeTechnique}
		}
		if i == 999 {
			produces = []reasoning.NodeType{reasoning.NodeTypeDataExposure}
		}
		classes = append(classes, reasoning.ActionClass{ID: id, Name: "synthetic", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: produces, RiskWeight: 0.01, ConfidenceBoost: 0.1})
	}
	eng.BindActionClasses(classes)

	start := time.Now()
	campaigns, err := eng.PlanCampaign(reasoning.NodeTypeDataExposure, reasoning.CampaignOptions{MaxDepth: 5, BeamWidth: 25, RiskTolerance: 3, ConfidenceThreshold: 0.2, TopN: 20})
	if err != nil {
		t.Fatalf("plan campaign: %v", err)
	}
	dur := time.Since(start)
	threshold := 5 * time.Second
	if dur > threshold {
		t.Fatalf("runtime exceeded threshold %v: %v", threshold, dur)
	}
	if len(campaigns) == 0 {
		t.Fatalf("expected at least one campaign")
	}
}
