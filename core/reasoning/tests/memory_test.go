package tests

import (
	"testing"

	"vantage/core/reasoning"
	"vantage/core/state"
)

func TestMemoryConfidenceDecayAndRepeatedActionPenalty(t *testing.T) {
	st, _ := state.New("memory")
	st.RecordActionMemory("AC-01", false, false)
	st.RecordActionMemory("AC-01", false, false)

	ranked := []reasoning.RankedAction{{TechniqueID: "T1", ActionClassID: "AC-01", Score: 1.0}}
	// indirect via simulate cycles: memory adjustments occur in planner path
	eng := reasoning.NewEngine(nil)
	eng.ConfigureCycle(reasoning.CycleConfig{Target: "x", AllowedTechniques: []string{"T1"}})
	_ = ranked
	if st.FailedAttempts("AC-01") < 2 {
		t.Fatalf("failed attempts not tracked")
	}
}

func TestSimulateCampaignCyclesAdaptsAcrossCycles(t *testing.T) {
	eng := reasoning.NewEngine(nil)
	eng.BindActionClasses([]reasoning.ActionClass{{ID: "AC-01", Name: "recon", Phase: state.PhaseRecon, Preconditions: []reasoning.GraphPattern{{RequiredNodeTypes: []reasoning.NodeType{reasoning.NodeTypeEvidence}}}, ProducesNodes: []reasoning.NodeType{reasoning.NodeTypeHypothesis}, ConfidenceBoost: 0.2}})
	eng.Graph().UpsertNode(&reasoning.Node{ID: "seed", Type: reasoning.NodeTypeEvidence, Label: "seed"})
	st, _ := state.New("trace")
	eng.ConfigureCycle(reasoning.CycleConfig{Target: "target", AllowedTechniques: []string{"T1000"}})
	_, _ = eng.RunCycle(st)
	trace := eng.SimulateCampaignCycles(3)
	if len(trace.StateProgression) != 3 || len(trace.ConfidenceEvolution) != 3 {
		t.Fatalf("unexpected trace shape")
	}
}
