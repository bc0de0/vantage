package reasoning

import "vantage/core/state"

type CampaignTrace struct {
	StateProgression    []state.Status
	PhaseTransitions    []state.OperationPhase
	ConfidenceEvolution []float64
}

func applyStateMemoryAdjustments(ranked []RankedAction, st *state.State) {
	if st == nil {
		return
	}
	for i := range ranked {
		fails := st.FailedAttempts(ranked[i].ActionClassID)
		if fails > 0 {
			ranked[i].Score -= float64(fails) * 0.2
		}
		if k, ok := st.ExposureKnowledge()[ranked[i].ActionClassID]; ok {
			ranked[i].Score += k * 0.1
		}
	}
}

func (e *Engine) SimulateCampaignCycles(n int) CampaignTrace {
	trace := CampaignTrace{StateProgression: make([]state.Status, 0, n), PhaseTransitions: make([]state.OperationPhase, 0, n), ConfidenceEvolution: make([]float64, 0, n)}
	if e == nil || n <= 0 {
		return trace
	}
	for i := 0; i < n; i++ {
		if e.state != nil {
			trace.StateProgression = append(trace.StateProgression, e.state.Status())
		} else {
			trace.StateProgression = append(trace.StateProgression, state.StatusInitialized)
		}
		decision, err := e.PlanNextAction(PlannerQuery{Target: e.cycle.Target, AllowedTechniques: e.cycle.AllowedTechniques, TopN: 1})
		if err != nil || decision == nil {
			trace.ConfidenceEvolution = append(trace.ConfidenceEvolution, 0)
			trace.PhaseTransitions = append(trace.PhaseTransitions, state.PhaseRecon)
			continue
		}
		trace.ConfidenceEvolution = append(trace.ConfidenceEvolution, decision.Selected.Score)
		trace.PhaseTransitions = append(trace.PhaseTransitions, state.PhaseRecon)
		if e.state != nil {
			e.state.RecordActionMemory(decision.Selected.ActionClassID, true, true)
		}
	}
	return trace
}
