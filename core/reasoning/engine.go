package reasoning

import (
	"fmt"

	"vantage/core/state"
)

// Planner ranks viable options using deterministic logic.
type Planner interface {
	Rank(phase state.OperationPhase, hypotheses []Hypothesis) []AttackOption
}

// OpsecEvaluator estimates detection pressure for candidate actions.
type OpsecEvaluator interface {
	Filter(options []AttackOption) []AttackOption
}

// Engine orchestrates one cognition-centric reasoning cycle.
type Engine struct {
	confidence *ConfidenceModel
	planner    Planner
	opsec      OpsecEvaluator
}

func NewEngine(planner Planner, opsec OpsecEvaluator) (*Engine, error) {
	if planner == nil {
		return nil, fmt.Errorf("planner is required")
	}
	if opsec == nil {
		return nil, fmt.Errorf("opsec evaluator is required")
	}
	return &Engine{
		confidence: NewConfidenceModel(),
		planner:    planner,
		opsec:      opsec,
	}, nil
}

func (e *Engine) RunCycle(input CycleInput) (CycleOutput, error) {
	if err := input.Phase.Validate(); err != nil {
		return CycleOutput{}, err
	}

	hypotheses := make([]Hypothesis, 0, len(input.Hypotheses))
	for _, hypothesis := range input.Hypotheses {
		hypothesis.Confidence = e.confidence.Score(
			hypothesis.SourceReliability,
			hypothesis.InferenceDepth,
			len(hypothesis.Evidence),
		)
		hypotheses = append(hypotheses, hypothesis)
	}

	options := e.planner.Rank(input.Phase, hypotheses)
	filtered := e.opsec.Filter(options)

	selectedPhase := input.Phase
	if len(filtered) == 0 {
		if next, ok := input.Phase.Next(); ok {
			selectedPhase = next
		}
	}

	return CycleOutput{
		SelectedPhase: selectedPhase,
		TopOptions:    filtered,
		Hypotheses:    hypotheses,
	}, nil
}
