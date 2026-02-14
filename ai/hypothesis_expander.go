package ai

import (
	"fmt"

	"vantage/core/reasoning"
	"vantage/core/state"
)

// HypothesisExpanderAdapter converts advisory AI output into reasoning hypotheses.
type HypothesisExpanderAdapter struct{}

// Expand returns advisory hypotheses from AI suggestions.
func (a HypothesisExpanderAdapter) Expand(graph *reasoning.Graph, _ *state.State) ([]reasoning.Hypothesis, error) {
	if graph == nil {
		return nil, nil
	}

	input := AdvisoryInput{}
	input.Intent.Objective = "Expand operational hypotheses from current evidence"
	input.Exposure.RemainingBudget = "unknown"
	input.TargetContext.HighLevelType = "unknown"

	evidenceNodes := graph.NodesByType(reasoning.NodeTypeEvidence)
	input.Intent.AllowedDomains = make([]string, 0, len(evidenceNodes))
	for _, n := range evidenceNodes {
		input.Intent.AllowedDomains = append(input.Intent.AllowedDomains, n.Label)
	}

	for id := range CanonicalActionClasses {
		input.CanonicalActionClasses = append(input.CanonicalActionClasses, id)
	}

	advisory, err := Advise(input)
	if err != nil {
		return nil, err
	}

	hypotheses := make([]reasoning.Hypothesis, 0, len(advisory.Suggested))
	for idx, suggestion := range advisory.Suggested {
		hypotheses = append(hypotheses, reasoning.Hypothesis{
			ID:                fmt.Sprintf("hyp-ai-%d-%s", idx+1, suggestion.ID),
			ActionClassID:     suggestion.ID,
			Statement:         fmt.Sprintf("AI advisory suggests %s: %s", suggestion.ID, suggestion.Rationale),
			SupportingNodeIDs: nil,
			Confidence:        suggestion.Confidence,
		})
	}

	return hypotheses, nil
}

var _ reasoning.HypothesisExpander = HypothesisExpanderAdapter{}
