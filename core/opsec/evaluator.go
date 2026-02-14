package opsec

import "vantage/core/reasoning"

// Evaluator removes noisy options likely to cross detection thresholds.
type Evaluator struct {
	minimumStealth float64
}

func NewEvaluator(minimumStealth float64) *Evaluator {
	return &Evaluator{minimumStealth: minimumStealth}
}

func (e *Evaluator) Filter(options []reasoning.AttackOption) []reasoning.AttackOption {
	filtered := make([]reasoning.AttackOption, 0, len(options))
	for _, option := range options {
		if option.StealthScore < e.minimumStealth {
			continue
		}
		filtered = append(filtered, option)
	}
	return filtered
}
