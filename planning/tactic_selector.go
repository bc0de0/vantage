package planning

import (
	"sort"

	"vantage/core/reasoning"
	"vantage/core/state"
)

// TacticSelector ranks options with weighted attack-path scoring.
type TacticSelector struct{}

func NewTacticSelector() *TacticSelector {
	return &TacticSelector{}
}

func (s *TacticSelector) Rank(phase state.OperationPhase, hypotheses []reasoning.Hypothesis) []reasoning.AttackOption {
	options := make([]reasoning.AttackOption, 0, len(hypotheses))
	for _, h := range hypotheses {
		impact := h.Confidence
		stealth := 1 - (0.12 * float64(h.InferenceDepth))
		if stealth < 0 {
			stealth = 0
		}
		effort := 0.5
		dependencies := 0.3
		score := (impact * 0.4) + (stealth * 0.3) + ((1 - effort) * 0.2) + ((1 - dependencies) * 0.1)

		options = append(options, reasoning.AttackOption{
			Name:           h.Claim,
			Phase:          phase,
			ImpactScore:    impact,
			StealthScore:   stealth,
			EffortScore:    effort,
			DependencyCost: dependencies,
			Score:          score,
			Rationale:      "Weighted impact/stealth/effort/dependency model",
		})
	}

	sort.Slice(options, func(i, j int) bool {
		return options[i].Score > options[j].Score
	})
	return options
}
