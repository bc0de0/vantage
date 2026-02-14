package reasoning

import (
	"fmt"
	"sort"
)

// Planner ranks action candidates using technique effects.
type Planner struct {
	registry TechniqueEffectRegistry
	weights  TechniqueScoreWeights
}

// NewPlanner creates a planner with a technique effect registry.
func NewPlanner(registry TechniqueEffectRegistry, weights TechniqueScoreWeights) *Planner {
	return &Planner{registry: registry, weights: weights}
}

// RankedActions returns sorted candidates for a target.
func (p *Planner) RankedActions(query PlannerQuery) []RankedAction {
	if p == nil || p.registry == nil {
		return nil
	}
	techniquesToScore := query.AllowedTechniques
	if len(techniquesToScore) == 0 {
		techniquesToScore = p.registry.KnownTechniques()
	}

	out := make([]RankedAction, 0, len(techniquesToScore))
	for _, id := range techniquesToScore {
		effect, ok := p.registry.EffectForTechnique(id)
		if !ok {
			continue
		}
		score := ScoreTechnique(effect, p.weights)
		out = append(out, RankedAction{
			TechniqueID: id,
			Target:      query.Target,
			Score:       score,
			Impact:      effect.Impact,
			Risk:        effect.Risk,
			Stealth:     effect.Stealth,
			Reason:      fmt.Sprintf("impact=%.2f risk=%.2f stealth=%.2f", effect.Impact, effect.Risk, effect.Stealth),
		})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Score == out[j].Score {
			return out[i].TechniqueID < out[j].TechniqueID
		}
		return out[i].Score > out[j].Score
	})

	if query.TopN > 0 && len(out) > query.TopN {
		out = out[:query.TopN]
	}
	return out
}
